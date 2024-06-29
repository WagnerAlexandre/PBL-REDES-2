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
	CPF1     string `json:"cpf1"`
	CPF2     string `json:"cpf2"`
	NumConta int    `json:"numconta"`
	Tipo     int    `json:"tipo"`
	Senha    string `json:"senha"`
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
	var conta ContaPF
	if err := c.BindJSON(&conta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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

	var novaConta *ContaPJ
	novaConta = criar_conta_pj(&cadastro)
	c.JSON(http.StatusCreated, novaConta)
}

func criar_conta_pj(novaConta *CadastroPFPJ) *ContaPJ {
	var conta ContaPJ
	conta.Nome = novaConta.Nome
	conta.CNPJ = novaConta.CPFCNPJ
	conta.Senha = novaConta.Senha

	conta.NumConta = criar_NumConta()
	conta.Balanco = 0.0

	// Armazena as informações da conta em TContasPJ
	TContasPJ[conta.NumConta] = &conta

	// Indexa a conta usando o CPF
	IndexcontasPJ[conta.CNPJ] = append(IndexcontasPJ[conta.CNPJ], &conta)
	return &conta
}

func rota_Cadastrar_CJ(c *gin.Context) {
	var conta ContaCJ
	if err := c.BindJSON(&conta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	criar_conta_cj(&conta)

	c.JSON(http.StatusCreated, gin.H{"message": "Conta conjunta criada com sucesso", "conta": conta})
}

func criar_conta_cj(conta *ContaCJ) {
	conta.NumConta = criar_NumConta()
	conta.Balanco = 0.0

	// Armazena as informações da conta em TContasCJ
	TContasCJ[conta.NumConta] = conta

	// Indexa a conta usando os CPFs
	IndexcontasCJ[conta.CPF1] = append(IndexcontasCJ[conta.CPF1], conta)
	IndexcontasCJ[conta.CPF2] = append(IndexcontasCJ[conta.CPF2], conta)
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
			if conta.Senha == req.Senha {
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
