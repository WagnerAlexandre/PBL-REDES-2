package main

// Desenolvido com o meu paradigma favorito: XGH

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

type Conta struct {
	ID      int     `json:"id"`
	Nome    string  `json:"nome"`
	Balanco float64 `json:"balanco"`
	Banco   string  `json:"Banco"`
}

// struct da operacao, se tipo == 1, soma o valor ao balanco, caso tipo == 2, reduz.
// ID - indica o numero da conta
// Valor -
type TransacaoWeb struct {
	ID    int     `json:"numConta"`
	Valor float64 `json:"valor"`
	Banco string  `json:"banco"`
	Tipo  int     `json:"tipo"`
}

// Execução do servidor

const (
	HOST   = "localhost"
	BBMN   = ":65500"
	BB     = ":65501"
	BG     = ":65502"
	ATUAL  = BG
	ENDPNT = "BG"
)

var bancoURLs = map[string]string{
	"BBMN": ":65500",
	"BB":   ":65501",
	"BG":   ":65502",
}

// Structs 2PC
type PrepareRequest struct {
	ID    int     `json:"id"`
	Valor float64 `json:"valor"`
	Tipo  int     `json:"tipo"`
	Url   string  `json:"url"`
}

type PrepareResponse struct {
	Status string `json:"status"`
}

type PreparedTransaction struct {
	ID    int
	Valor float64
	Tipo  int
	Url   string
}

var preparedTransactions = make(map[int]PreparedTransaction)
var preparedMutex sync.Mutex

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

	// rota para a página de transferência
	router.GET("/transferencia", func(c *gin.Context) {
		c.File("./static/transferencia.html")
	})

	// Rotas da API
	router.POST("/login", loginHandler)
	router.POST("/criarContaPF", rota_Cadastrar_PF)
	router.POST("/criarContaPJ", rota_Cadastrar_PJ)
	router.POST("/criarContaCJ", rota_Cadastrar_CJ)
	router.POST("/somaLocal", somaLocalHandler)
	router.POST("/reducaoLocal", reducaoLocalHandler)
	router.POST("/getContas", retornarContas)
	router.POST("/procurarConta", procuraConta)
	router.POST("/realizarTransferencia", iniciarTransferencia)

	//Rotas 2PC
	router.POST("/prepare", prepareHandler)
	router.POST("/commit", commitHandler)
	router.POST("/abort", abortHandler)

	// Manipulação admin
	router.GET("/contasPF", getContasPF)
	router.GET("/contasPJ", getContasPJ)
	router.GET("/contasCJ", getContasCJ)

	// ROTA Query, para procurar contas relacionadas a um cpf/cnpj
	// Rota para retornar contas associadas ao CPF/CNPJ do cliente
	router.GET("/getUnicaChaveContas", func(c *gin.Context) {
		cpfCnpj := c.Query("cpf_cnpj")

		var contas []Conta
		// Buscar as contas PF associadas ao CPF
		// Buscar as contas PF associadas ao CPF
		if contasPF, exists := IndexcontasPF[cpfCnpj]; exists {
			for _, conta := range contasPF {
				contas = append(contas, Conta{
					ID:      conta.NumConta,
					Nome:    conta.Nome,
					Balanco: conta.Balanco,
					Banco:   ENDPNT,
				})
			}
		}

		// Buscar as contas CJ associadas ao CPF
		if contasCJ, exists := IndexcontasCJ[cpfCnpj]; exists {
			for _, conta := range contasCJ {
				contas = append(contas, Conta{
					ID:      conta.NumConta,
					Nome:    conta.Nome,
					Balanco: conta.Balanco,
					Banco:   ENDPNT,
				})
			}
		}

		// Buscar as contas PJ associadas ao CNPJ
		if contasPJ, exists := IndexcontasPJ[cpfCnpj]; exists {
			for _, conta := range contasPJ {
				contas = append(contas, Conta{
					ID:      conta.NumConta,
					Nome:    conta.Nome,
					Balanco: conta.Balanco,
					Banco:   ENDPNT,
				})
			}
		}
		c.JSON(http.StatusOK, contas)
	})

	fmt.Printf("Servidor rodando em http://%s:%s\n", HOST, ATUAL)
	if err := router.Run(HOST + ATUAL); err != nil {
		panic(err)
	}
}

// funcao de ajuda para converter da estrutura prepare request para json
func (r PrepareRequest) toJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

// Prepare handler
func prepareHandler(c *gin.Context) {
	var request PrepareRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar json ao preparar commit"})
		return
	}

	var conta interface{}
	var exists bool

	conta, exists = TContasPF[request.ID]
	if !exists {
		conta, exists = TContasPJ[request.ID]
		if !exists {
			conta, exists = TContasCJ[request.ID]
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
				return
			}
		}
	}

	var balanco *float64
	var mutex *sync.Mutex

	switch v := conta.(type) {
	case *ContaPF:
		balanco = &v.Balanco
		mutex = &v.mutex
	case *ContaPJ:
		balanco = &v.Balanco
		mutex = &v.mutex
	case *ContaCJ:
		balanco = &v.Balanco
		mutex = &v.mutex
	}

	mutex.Lock()
	defer mutex.Unlock()

	if request.Tipo == 2 { // Redução de saldo
		if *balanco < request.Valor {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente"})
			return
		}
	}

	preparedMutex.Lock()
	preparedTransactions[request.ID] = PreparedTransaction{
		ID:    request.ID,
		Valor: request.Valor,
		Tipo:  request.Tipo,
		Url:   request.Url,
	}
	preparedMutex.Unlock()

	c.JSON(http.StatusOK, PrepareResponse{Status: "OK"})
}

// Commit handler
func commitHandler(c *gin.Context) {
	var request PrepareRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar json na fase de commit"})
		return
	}

	preparedMutex.Lock()
	transaction, exists := preparedTransactions[request.ID]
	if !exists {
		preparedMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	delete(preparedTransactions, request.ID)
	preparedMutex.Unlock()

	var conta interface{}
	var balanco *float64
	var mutex *sync.Mutex

	conta, exists = TContasPF[transaction.ID]
	if !exists {
		conta, exists = TContasPJ[transaction.ID]
		if !exists {
			conta, exists = TContasCJ[transaction.ID]
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada na fase de commit"})
				return
			}
		}
	}

	switch v := conta.(type) {
	case *ContaPF:
		balanco = &v.Balanco
		mutex = &v.mutex
	case *ContaPJ:
		balanco = &v.Balanco
		mutex = &v.mutex
	case *ContaCJ:
		balanco = &v.Balanco
		mutex = &v.mutex
	}

	mutex.Lock()
	defer mutex.Unlock()

	if transaction.Tipo == 2 { // Redução de saldo
		*balanco -= transaction.Valor
	} else if transaction.Tipo == 1 { // Aumento de saldo
		*balanco += transaction.Valor
	}

	c.JSON(http.StatusOK, gin.H{"status": "committed"})
}

// Abort handler
func abortHandler(c *gin.Context) {
	var request PrepareRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar json de conta unica"})
		return
	}

	preparedMutex.Lock()
	// se a transacao preparada ainda existe
	_, exists := preparedTransactions[request.ID]
	if !exists {
		preparedMutex.Unlock()
		// If the transaction does not exist, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	// Remove a transacao das transacoes preparadas
	delete(preparedTransactions, request.ID)
	// desbloqueia o mutex
	preparedMutex.Unlock()

	// retorna sucesso
	c.JSON(http.StatusOK, gin.H{"status": "aborted"})
}

// Funções auxiliares para realizar as transações com 2PC
func prepararTransacao(url string, request PrepareRequest) (bool, error) {
	resp, err := http.Post(fmt.Sprintf("http://localhost%s/prepare", url), "application/json", bytes.NewBuffer(request.toJSON()))
	if err != nil {

		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("prepare request failed")
	}
	return true, nil
}

func commitTransacao(url string, request PrepareRequest) error {
	resp, err := http.Post(fmt.Sprintf("http://localhost%s/commit", url), "application/json", bytes.NewBuffer(request.toJSON()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("commit request failed")
	}
	return nil
}

func abortTransacao(url string, request PrepareRequest) error {
	resp, err := http.Post(fmt.Sprintf("http://localhost%s/abort", url), "application/json", bytes.NewBuffer(request.toJSON()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("abort request failed")
	}
	return nil
}

// Função para iniciar a transferência utilizando 2PC
func iniciarTransferencia(c *gin.Context) {
	var requests []TransacaoWeb

	if err := c.ShouldBindJSON(&requests); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar json pacote de contas"})
		return
	}

	prepares := []PrepareRequest{}

	for _, request := range requests {

		prepares = append(prepares, PrepareRequest{
			ID:    request.ID,
			Valor: request.Valor,
			Tipo:  request.Tipo,
			Url:   bancoURLs[request.Banco],
		})

	}

	allPrepared := true
	for _, prepare := range prepares {
		ok, err := prepararTransacao(prepare.Url, prepare)
		if err != nil || !ok {
			allPrepared = false
			break
		}
	}

	if allPrepared {
		for _, prepare := range prepares {
			if err := commitTransacao(prepare.Url, prepare); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Commit failed"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "Transferências realizadas com sucesso"})
	} else {
		for _, prepare := range prepares {
			_ = abortTransacao(prepare.Url, prepare)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transferências abortadas"})
	}
}

// retorna contas pertencentes ao cpf/cnpj do banco BBMN

func getUnicaChaveContasFromBBMN(cpfCnpj string) ([]Conta, error) {
	resp, err := http.Get("http://localhost:65500/getUnicaChaveContas?cpf_cnpj=" + cpfCnpj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contas []Conta
	err = json.NewDecoder(resp.Body).Decode(&contas)
	if err != nil {
		return nil, err
	}

	return contas, nil
}

// retorna contas pertencentes ao cpf/cnpj do banco BB
func getUnicaChaveContasFromBB(cpfCnpj string) ([]Conta, error) {
	resp, err := http.Get("http://localhost:65501/getUnicaChaveContas?cpf_cnpj=" + cpfCnpj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contas []Conta
	err = json.NewDecoder(resp.Body).Decode(&contas)
	if err != nil {
		return nil, err
	}

	return contas, nil
}

// formato de requisicao das informacoes de uma determinada conta
type ContaRequest struct {
	NumConta int    `json:"numconta"`
	Banco    string `json:"banco"`
}

type ContaResponse struct {
	NumConta int    `json:"numconta"`
	Nome     string `json:"nome"`
	Banco    string `json:"banco"`
}

// procura a conta "numconta" do Banco "banco" para ser o alvo da transferencia.

func procuraConta(c *gin.Context) {
	var request ContaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	banco := request.Banco
	numConta := request.NumConta
	if banco == ENDPNT {
		fd := 0

		// Procura localmente
		if contaPF, ok := TContasPF[numConta]; ok {
			cookieValue := fmt.Sprintf("%s|%d", ENDPNT, contaPF.NumConta)

			c.SetCookie("alvoTransacao", cookieValue, 3600, "/", "localhost", false, false)

			c.JSON(http.StatusOK, ContaResponse{NumConta: contaPF.NumConta, Nome: contaPF.Nome, Banco: ENDPNT})
			return
		}
		if contaPJ, ok := TContasPJ[numConta]; ok {
			cookieValue := fmt.Sprintf("%s|%d", ENDPNT, contaPJ.NumConta)

			c.SetCookie("alvoTransacao", cookieValue, 3600, "/", "localhost", false, false)
			c.JSON(http.StatusOK, ContaResponse{NumConta: contaPJ.NumConta, Nome: contaPJ.Nome, Banco: ENDPNT})

			return

		}
		if contaCJ, ok := TContasCJ[numConta]; ok {
			cookieValue := fmt.Sprintf("%s|%d", ENDPNT, contaCJ.NumConta)

			c.SetCookie("alvoTransacao", cookieValue, 3600, "/", "localhost", false, false)
			c.JSON(http.StatusOK, ContaResponse{NumConta: contaCJ.NumConta, Nome: contaCJ.Nome, Banco: ENDPNT})

			return
		}
		println(fd)
		c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		return
	}

	// Redireciona para outros servidores
	var endpoint string
	if banco == "BBMN" {
		endpoint = fmt.Sprintf("http://%s%s/procurarConta", HOST, BBMN)
	} else if banco == "BB" {
		endpoint = fmt.Sprintf("http://%s%s/procurarConta", HOST, BB)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Banco desconhecido"})
		return
	}

	client := &http.Client{}
	jsonData, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar JSON"})
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar requisição"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer requisição"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler resposta"})
		return
	}

	if resp.StatusCode != http.StatusOK {

		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	var response ContaResponse
	if err := json.Unmarshal(body, &response); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao analisar resposta"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// retorna as contas associadas a um determinado cpf/cnjp logado

func retornarContas(c *gin.Context) {
	// Obter o valor do cookie
	cookie, err := c.Cookie("gringuesco")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não foi possível obter o cookie"})
		return
	}

	// Extrair o CPF/CNPJ do cookie
	cookieParts := strings.Split(cookie, "|")
	if len(cookieParts) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato do cookie inválido"})
		return
	}
	cpfCnpj := cookieParts[1]

	// Inicializar o slice de contas
	var contas []Conta

	// Buscar as contas PF associadas ao CPF
	if contasPF, exists := IndexcontasPF[cpfCnpj]; exists {
		for _, conta := range contasPF {
			contas = append(contas, Conta{
				ID:      conta.NumConta,
				Nome:    conta.Nome,
				Balanco: conta.Balanco,
				Banco:   ENDPNT,
			})
		}
	}

	// Buscar as contas CJ associadas ao CPF
	if contasCJ, exists := IndexcontasCJ[cpfCnpj]; exists {
		for _, conta := range contasCJ {
			contas = append(contas, Conta{
				ID:      conta.NumConta,
				Nome:    conta.Nome,
				Balanco: conta.Balanco,
				Banco:   ENDPNT,
			})
		}
	}

	// Buscar as contas PJ associadas ao CNPJ
	if contasPJ, exists := IndexcontasPJ[cpfCnpj]; exists {
		for _, conta := range contasPJ {
			contas = append(contas, Conta{
				ID:      conta.NumConta,
				Nome:    conta.Nome,
				Balanco: conta.Balanco,
				Banco:   ENDPNT,
			})
		}
	}

	// Conectar ao banco BBMN e adicionar suas contas
	bbmnContas, err := getUnicaChaveContasFromBBMN(cpfCnpj)
	if err == nil {
		contas = append(contas, bbmnContas...)
	}

	// Conectar ao banco BG e adicionar suas contas
	bgContas, err := getUnicaChaveContasFromBB(cpfCnpj)
	if err == nil {
		contas = append(contas, bgContas...)
	}

	// Verificar se há contas
	if len(contas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma conta encontrada"})
		return
	}

	// Retornar as contas em formato JSON
	c.JSON(http.StatusOK, contas)
}

// altera saldo local

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

			if (contaPJ.Balanco - transacao.Valor) >= 0 {
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

// Funcoes para criacao de contas

func criar_NumConta() int {
	mutexCriacao.Lock()
	defer mutexCriacao.Unlock()
	temp := numero_contas
	numero_contas++
	return temp
}

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
		c.SetCookie("gringuesco", cookieValue, 3600, "/", "localhost", false, false)

		c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
	}
}

// checa se o usuario já possui conta com cpf ou cnpj no banco atual

func checaExistencia(cpfcnpj string) bool {
	_, existsPF := IndexcontasPF[cpfcnpj]
	_, existsPJ := IndexcontasPJ[cpfcnpj]
	return existsPF || existsPJ
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
