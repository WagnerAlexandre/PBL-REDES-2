package main

import (
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Estrutura para conta bancaria pessoa fisica
type ContaPF struct {
	NumConta int     `json:"numconta"`
	CPF      string  `json:"cpf"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
}

// Estrutura para conta bancaria pessoa juridica
type ContaPJ struct {
	NumConta int     `json:"numconta"`
	CNPJ     string  `json:"cnpj"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
}

// Estrutura para conta bancaria conjunta
type ContaCJ struct {
	NumConta int     `json:"numconta"`
	CPF1     string  `json:"cpf1"`
	CPF2     string  `json:"cpf2"`
	Nome     string  `json:"nome"`
	Senha    string  `json:"senha"`
	Tipo     int     `json:"tipo"`
	Balanco  float64 `json:"balanco"`
}

// Estrutura representando o login
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

// mutex para numero de conta
var mutexCriacao = &sync.Mutex{}
var numero_contas int = 1

// Função para criar um número de conta único de forma segura
func criar_NumConta() int {
	mutexCriacao.Lock()         // Bloqueia o mutex antes de modificar numero_contas
	defer mutexCriacao.Unlock() // Garante que o mutex seja desbloqueado após o término da função
	temp := numero_contas
	numero_contas++
	return temp
}

var TContasPF = make(map[int]*ContaPF)
var TContasPJ = make(map[int]*ContaPJ)
var TContasCJ = make(map[int]*ContaCJ)

var IndexcontasPF = make(map[string][]*ContaPF)
var IndexcontasPJ = make(map[string][]*ContaPJ)
var IndexcontasCJ = make(map[string][]*ContaCJ)

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
	c.JSON(http.StatusCreated, gin.H{"message": "Conta criada com sucesso", "conta": conta})
}

func criar_conta_pf(conta *ContaPF) {
	conta.NumConta = criar_NumConta()
	conta.Balanco = 0.0

	// Armazena as informações da conta em TContasPF
	TContasPF[conta.NumConta] = conta

	// Indexa a conta usando o CPF
	IndexcontasPF[conta.CPF] = append(IndexcontasPF[conta.CPF], conta)
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
	c.JSON(http.StatusCreated, novaConta)
}

func criar_conta_pj(cadastro *CadastroPFPJ) *ContaPJ {
	conta := &ContaPJ{
		NumConta: criar_NumConta(),
		CNPJ:     cadastro.CPFCNPJ,
		Nome:     cadastro.Nome,
		Senha:    cadastro.Senha,
		Tipo:     2, // Tipo para Pessoa Jurídica
		Balanco:  0.0,
	}

	// Armazena as informações da conta em TContasPJ
	TContasPJ[conta.NumConta] = conta

	// Indexa a conta usando o CNPJ
	IndexcontasPJ[conta.CNPJ] = append(IndexcontasPJ[conta.CNPJ], conta)

	return conta
}

func criar_conta_cj(cadastro *CadastroCJ) *ContaCJ {
	conta := &ContaCJ{
		NumConta: criar_NumConta(),
		CPF1:     cadastro.CPF1,
		CPF2:     cadastro.CPF2,
		Nome:     cadastro.CPF1 + " e " + cadastro.CPF2, // Opcional
		Senha:    cadastro.Senha,
		Tipo:     3, // Tipo para Conta Conjunta
		Balanco:  0.0,
	}

	// Armazena as informações da conta em TContasCJ
	TContasCJ[conta.NumConta] = conta

	// Indexa a conta usando os CPFs
	IndexcontasCJ[conta.CPF1] = append(IndexcontasCJ[conta.CPF1], conta)
	IndexcontasCJ[conta.CPF2] = append(IndexcontasCJ[conta.CPF2], conta)

	return conta
}

func rota_Cadastrar_CJ(c *gin.Context) {
	var cadastro CadastroCJ
	if err := c.BindJSON(&cadastro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaConta := criar_conta_cj(&cadastro)
	c.JSON(http.StatusCreated, gin.H{"message": "Conta conjunta criada com sucesso", "conta": novaConta})
}

func checaExistencia(chave string) bool {
	// Verifica se a chave existe em IndexcontasPF
	_, existsPF := IndexcontasPF[chave]
	// Verifica se a chave existe em IndexcontasPJ
	_, existsPJ := IndexcontasPJ[chave]

	// Retorna true se a chave existe em IndexcontasPF ou IndexcontasPJ
	return existsPF || existsPJ
}

func main() {
	router := gin.Default()
	// Configuração básica do CORS para permitir todas as origens
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"} // Incluir OPTIONS para tratamento de CORS

	router.Use(cors.New(config))
	// ROTAS
	router.GET("/contasPF", getContasPF)
	router.GET("/contasPJ", getContasPJ)
	router.GET("/contasCJ", getContasCJ)
	router.POST("/login", loginHandler)
	router.POST("/criarContaPF", rota_Cadastrar_PF)
	router.POST("/criarContaPJ", rota_Cadastrar_PJ)
	router.POST("/criarContaCJ", rota_Cadastrar_CJ)

	router.Run("localhost:8080")
}

func loginHandler(c *gin.Context) {
	var req Login
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Verifica o tipo de conta e seleciona o mapa apropriado

	switch req.Tipo {
	case 1:
		if conta, exists := TContasPF[req.NumConta]; exists {
			if conta.Senha == req.Senha {
				c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido"})
				return
			}
		}
	case 2:
		if conta, exists := TContasPJ[req.NumConta]; exists {
			if conta.Senha == req.Senha && conta.CNPJ == req.CPFRAZAO {
				c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido"})
				return
			}
		}
	case 3:
		if conta, exists := TContasCJ[req.NumConta]; exists {
			if conta.Senha == req.Senha {
				c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido"})
				return
			}
		}
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de conta inválido"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
}

func getContasPF(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, IndexcontasPF)
}

func getContasPJ(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, IndexcontasPJ)
}

func getContasCJ(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, IndexcontasCJ)
}
