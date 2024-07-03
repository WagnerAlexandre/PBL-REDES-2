package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ContaPF struct {
	NumConta int     `json:"numconta"`
	CPF      string  `json:"cpf"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
	mutex    sync.Mutex
}

type ContaPJ struct {
	NumConta int     `json:"numconta"`
	CNPJ     string  `json:"cnpj"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
	mutex    sync.Mutex
}

type ContaCJ struct {
	NumConta int     `json:"numconta"`
	CPF1     string  `json:"cpf1"`
	CPF2     string  `json:"cpf2"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
	mutex    sync.Mutex
}

type Login struct {
	CPFRAZAO string `json:"CPFRAZAO"`
	NumConta int    `json:"numconta"`
	Tipo     int    `json:"tipo"`
	Senha    string `json:"senha"`
}

type CadastroPFPJ struct {
	CPFCNPJ string `json:"cpfcnpj"`
	Nome    string `json:"nome"`
	Senha   string `json:"senha"`
}

type CadastroCJ struct {
	CPF1  string `json:"cpf1"`
	CPF2  string `json:"cpf2"`
	Senha string `json:"senha"`
}

type Transacao struct {
	IDConta   int     `json:"idConta"`
	TipoConta int     `json:"tipoconta"`
	Valor     float64 `json:"valor"`
}

var (
	mutexCriacao  = &sync.Mutex{}
	numero_contas = 1
)

var (
	TContasPF     = make(map[int]*ContaPF)
	TContasPJ     = make(map[int]*ContaPJ)
	TContasCJ     = make(map[int]*ContaCJ)
	IndexcontasPF = make(map[string][]*ContaPF)
	IndexcontasPJ = make(map[string][]*ContaPJ)
	IndexcontasCJ = make(map[string][]*ContaCJ)
)

func main() {
	router := gin.Default()

	// Configuração básica do CORS para permitir todas as origens
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"} // Incluir OPTIONS para tratamento de CORS
	router.Use(cors.New(config))

	// Servir arquivos estáticos
	router.Static("/static", "./static")

	// Endpoints para páginas HTML
	router.GET("/inicial", func(c *gin.Context) {
		c.File("./static/login.html")
	})

	router.GET("/cadastro", func(c *gin.Context) {
		c.File("./static/cadastro.html")
	})

	router.GET("/menuprincipal", func(c *gin.Context) {
		c.File("./static/menuprincipal.html")
	})

	// Rotas da API
	router.POST("/login", loginHandler)
	router.POST("/criarContaPF", rota_Cadastrar_PF)
	router.POST("/criarContaPJ", rota_Cadastrar_PJ)
	router.POST("/criarContaCJ", rota_Cadastrar_CJ)
	router.POST("/somaLocal", somaLocalHandler)
	router.POST("/reducaoLocal", reducaoLocalHandler)

	// Manipulação admin
	router.GET("/contasPF", getContasPF)
	router.GET("/contasPJ", getContasPJ)
	router.GET("/contasCJ", getContasCJ)

	// Execução do servidor
	serverPort := ":65501"
	const (
		HOST  = "0.0.0.0"
		BBMN  = "65500"
		BB    = "65501"
		BG    = "65502"
		ATUAL = BB
	)
	fmt.Printf("Servidor rodando em http://%s%s\n", HOST, serverPort)
	if err := router.Run(HOST + serverPort); err != nil {
		panic(err)
	}
}

func reducaoLocalHandler(c *gin.Context) {
	var transacao Transacao
	if err := c.BindJSON(&transacao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if transacao.TipoConta == 1 {
		if contaPF, exists := TContasPF[transacao.IDConta]; exists {
			contaPF.mutex.Lock()
			defer contaPF.mutex.Unlock()

			if (contaPF.Balanco - transacao.Valor) >= 0 {
				contaPF.Balanco -= transacao.Valor
				c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaPF.Balanco})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Saldo insuficiente"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	} else if transacao.TipoConta == 2 {
		if contaPJ, exists := TContasPJ[transacao.IDConta]; exists {
			contaPJ.mutex.Lock()
			defer contaPJ.mutex.Unlock()

			if contaPJ.Balanco >= transacao.Valor {
				contaPJ.Balanco -= transacao.Valor
				c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaPJ.Balanco})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Saldo insuficiente"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	} else if transacao.TipoConta == 3 {
		if contaCJ, exists := TContasCJ[transacao.IDConta]; exists {
			contaCJ.mutex.Lock()
			defer contaCJ.mutex.Unlock()

			if contaCJ.Balanco >= transacao.Valor {
				contaCJ.Balanco -= transacao.Valor
				c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaCJ.Balanco})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Saldo insuficiente"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	}
}

func somaLocalHandler(c *gin.Context) {
	var transacao Transacao
	if err := c.BindJSON(&transacao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if transacao.TipoConta == 1 {
		if contaPF, exists := TContasPF[transacao.IDConta]; exists {
			contaPF.mutex.Lock()
			defer contaPF.mutex.Unlock()

			contaPF.Balanco += transacao.Valor
			c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaPF.Balanco})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	} else if transacao.TipoConta == 2 {
		if contaPJ, exists := TContasPJ[transacao.IDConta]; exists {
			contaPJ.mutex.Lock()
			defer contaPJ.mutex.Unlock()

			contaPJ.Balanco += transacao.Valor
			c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaPJ.Balanco})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	} else if transacao.TipoConta == 3 {
		if contaCJ, exists := TContasCJ[transacao.IDConta]; exists {
			contaCJ.mutex.Lock()
			defer contaCJ.mutex.Unlock()

			contaCJ.Balanco += transacao.Valor
			c.JSON(http.StatusOK, gin.H{"success": true, "newBalance": contaCJ.Balanco})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		}
	}
}

func criar_NumConta() int {
	mutexCriacao.Lock()
	defer mutexCriacao.Unlock()
	temp := numero_contas
	numero_contas++
	return temp
}

// Funcoes para criacao de contas

func criar_conta_pf(conta *ContaPF) {
	conta.NumConta = criar_NumConta()
	conta.Balanco = 0.0
	TContasPF[conta.NumConta] = conta
	IndexcontasPF[conta.CPF] = append(IndexcontasPF[conta.CPF], conta)
}

func criar_conta_pj(cadastro *CadastroPFPJ) *ContaPJ {
	conta := &ContaPJ{
		NumConta: criar_NumConta(),
		CNPJ:     cadastro.CPFCNPJ,
		Nome:     cadastro.Nome,
		Senha:    cadastro.Senha,
		Tipo:     2,
		Balanco:  0.0,
	}
	TContasPJ[conta.NumConta] = conta
	IndexcontasPJ[conta.CNPJ] = append(IndexcontasPJ[conta.CNPJ], conta)
	return conta
}

func criar_conta_cj(cadastro *CadastroCJ) *ContaCJ {
	conta := &ContaCJ{
		NumConta: criar_NumConta(),
		CPF1:     cadastro.CPF1,
		CPF2:     cadastro.CPF2,
		Nome:     cadastro.CPF1 + " e " + cadastro.CPF2,
		Senha:    cadastro.Senha,
		Tipo:     3,
		Balanco:  0.0,
	}
	TContasCJ[conta.NumConta] = conta
	IndexcontasCJ[conta.CPF1] = append(IndexcontasCJ[conta.CPF1], conta)
	IndexcontasCJ[conta.CPF2] = append(IndexcontasCJ[conta.CPF2], conta)
	return conta
}

func rota_Cadastrar_PF(c *gin.Context) {
	var cadastro CadastroPFPJ
	if err := c.BindJSON(&cadastro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if checaExistencia(cadastro.CPFCNPJ) {
		c.JSON(http.StatusConflict, gin.H{"error": "Você já possui uma conta com este CPF"})
		return
	}

	var conta ContaPF
	conta.CPF = cadastro.CPFCNPJ
	conta.Nome = cadastro.Nome
	conta.Senha = cadastro.Senha
	conta.Tipo = 1 // Tipo para Pessoa Física

	criar_conta_pf(&conta)

	// Converter o número da conta para string antes de salvar o cookie
	numeroContaStr := strconv.Itoa(conta.NumConta)

	// Salvando cookie com o número da conta
	salvarCookieNumConta(c, numeroContaStr)

	c.JSON(http.StatusCreated, gin.H{"message": "Conta criada com sucesso"})
}

func rota_Cadastrar_PJ(c *gin.Context) {
	var cadastro CadastroPFPJ
	if err := c.BindJSON(&cadastro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if checaExistencia(cadastro.CPFCNPJ) {
		c.JSON(http.StatusConflict, gin.H{"error": "Você já possui uma conta com este CNPJ"})
		return
	}

	novaConta := criar_conta_pj(&cadastro)

	// Converter o número da conta para string antes de salvar o cookie
	numeroContaStr := strconv.Itoa(novaConta.NumConta)

	// Salvando cookie com o número da conta
	salvarCookieNumConta(c, numeroContaStr)

	c.JSON(http.StatusCreated, novaConta)
}

func rota_Cadastrar_CJ(c *gin.Context) {
	var cadastro CadastroCJ
	if err := c.BindJSON(&cadastro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaConta := criar_conta_cj(&cadastro)

	// Converter o número da conta para string antes de salvar o cookie
	numeroContaStr := strconv.Itoa(novaConta.NumConta)

	// Salvando cookie com o número da conta
	salvarCookieNumConta(c, numeroContaStr)

	c.JSON(http.StatusCreated, gin.H{"message": "Conta conjunta criada com sucesso", "conta": novaConta})
}

func salvarCookieNumConta(c *gin.Context, NumConta string) {
	cookie := &http.Cookie{
		Name:     "NumConta",
		Value:    NumConta,
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(c.Writer, cookie)
}

func loginHandler(c *gin.Context) {
	var req Login
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var conta interface{}
	var nome, cpfRazao string
	var balanco float64
	var tipo int

	switch req.Tipo {
	case 1:
		if c, exists := TContasPF[req.NumConta]; exists && c.Senha == req.Senha {
			conta = c
			nome = c.Nome
			cpfRazao = c.CPF
			balanco = c.Balanco
			tipo = 1
		}
	case 2:
		if c, exists := TContasPJ[req.NumConta]; exists && c.Senha == req.Senha && c.CNPJ == req.CPFRAZAO {
			conta = c
			nome = c.Nome
			cpfRazao = c.CNPJ
			balanco = c.Balanco
			tipo = 2
		}
	case 3:
		if c, exists := TContasCJ[req.NumConta]; exists && c.Senha == req.Senha {
			conta = c
			nome = c.Nome
			cpfRazao = c.CPF1 + " e " + c.CPF2
			balanco = c.Balanco
			tipo = 3
		}
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de conta inválido"})
		return
	}

	if conta != nil {
		// Construir o valor do cookie
		cookieValue := fmt.Sprintf("%s|%s|%d|%.2f|%d", nome, cpfRazao, req.NumConta, balanco, tipo)
		// Definir o cookie no contexto da requisição
		c.SetCookie("brasilheirinho", cookieValue, 3600, "/", "localhost", false, false)

		c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
	}
}

func getContasPF(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contasPF": TContasPF})
}

func getContasPJ(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contasPJ": TContasPJ})
}

func getContasCJ(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contasCJ": TContasCJ})
}

func checaExistencia(cpfcnpj string) bool {
	_, existsPF := IndexcontasPF[cpfcnpj]
	_, existsPJ := IndexcontasPJ[cpfcnpj]
	return existsPF || existsPJ
}
