package main

type usuario struct {
	CPF       string  `json:"cpf"`
	Nome      string  `json:"nome"`
	Senha     string  `json:"senha"`
	Email     string  `json:"email"`
	Balanco   float64 `json:"balanco"`
	Semaforo  bool    `json:"semaforo"`
	TipoConta int     `json:"tipoconta"`
}

var usuarios = []usuario{
	{CPF: "1", Nome: "Blue Train", Senha: "John Coltrane", Email: "inv@inv.cs", Balanco: 56.99, Semaforo: false},
	{CPF: "2", Nome: "Blue Train", Senha: "John Coltrane", Email: "inv@inv.cs", Balanco: 56.99, Semaforo: false},
	{CPF: "3", Nome: "Blue Train", Senha: "John Coltrane", Email: "inv@inv.cs", Balanco: 56.99, Semaforo: false},
	{CPF: "4", Nome: "Blue Train", Senha: "John Coltrane", Email: "inv@inv.cs", Balanco: 56.99, Semaforo: false},
}
