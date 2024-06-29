#!/bin/bash

# URL base do servidor
BASE_URL="http://localhost:8080"

# Dados de exemplo para contas PF, PJ e CJ
ACCOUNTS_PF=(
  '{"cpf": "11111111111", "nome": "Pessoa 1", "senha": "senha1", "numconta": 0, "tipo": 1, "balanco": 100.0}'
  '{"cpf": "22222222222", "nome": "Pessoa 2", "senha": "senha2", "numconta": 0, "tipo": 1, "balanco": 200.0}'
  '{"cpf": "33333333333", "nome": "Pessoa 3", "senha": "senha3", "numconta": 0, "tipo": 1, "balanco": 300.0}'
  '{"cpf": "44444444444", "nome": "Pessoa 4", "senha": "senha4", "numconta": 0, "tipo": 1, "balanco": 400.0}'
  '{"cpf": "55555555555", "nome": "Pessoa 5", "senha": "senha5", "numconta": 0, "tipo": 1, "balanco": 500.0}'
  '{"cpf": "66666666666", "nome": "Pessoa 6", "senha": "senha6", "numconta": 0, "tipo": 1, "balanco": 600.0}'
)

ACCOUNTS_PJ=(
  '{"cnpj": "11111111000100", "nome": "Empresa 1", "senha": "senha1", "numconta": 0, "tipo": 2, "balanco": 1000.0}'
  '{"cnpj": "22222222000100", "nome": "Empresa 2", "senha": "senha2", "numconta": 0, "tipo": 2, "balanco": 2000.0}'
  '{"cnpj": "33333333000100", "nome": "Empresa 3", "senha": "senha3", "numconta": 0, "tipo": 2, "balanco": 3000.0}'
  '{"cnpj": "44444444000100", "nome": "Empresa 4", "senha": "senha4", "numconta": 0, "tipo": 2, "balanco": 4000.0}'
  '{"cnpj": "55555555000100", "nome": "Empresa 5", "senha": "senha5", "numconta": 0, "tipo": 2, "balanco": 5000.0}'
  '{"cnpj": "66666666000100", "nome": "Empresa 6", "senha": "senha6", "numconta": 0, "tipo": 2, "balanco": 6000.0}'
)

# Função para criar contas PF sequencialmente
create_accounts_pf() {
  for account in "${ACCOUNTS_PF[@]}"; do
    curl -s -X POST "$BASE_URL/criarContaPF" -H "Content-Type: application/json" -d "$account"
    echo "" # Apenas para separar as respostas das requisições no output
  done
}

# Função para criar contas PJ sequencialmente
create_accounts_pj() {
  for account in "${ACCOUNTS_PJ[@]}"; do
    curl -s -X POST "$BASE_URL/criarContaPJ" -H "Content-Type: application/json" -d "$account"
    echo "" # Apenas para separar as respostas das requisições no output
  done
}

# Executa as funções de criação de contas sequencialmente
create_accounts_pf
create_accounts_pj

echo "Testes de criação de contas PF e PJ finalizados sequencialmente."
