#!/bin/bash

# Função para criar conta PF sequencialmente
create_account_pf() {
  cpf="$1"
  nome="$2"
  senha="$3"
  curl -s -X POST \
    http://localhost:65501/criarContaPF \
    -H 'Content-Type: application/json' \
    -d '{
      "cpfcnpj": "'"$cpf"'",
      "nome": "'"$nome"'",
      "senha": "'"$senha"'"
    }'
}

# Função para criar conta PJ sequencialmente
create_account_pj() {
  cnpj="$1"
  nome="$2"
  senha="$3"
  curl -s -X POST \
    http://localhost:65501/criarContaPJ \
    -H 'Content-Type: application/json' \
    -d '{
      "cpfcnpj": "'"$cnpj"'",
      "nome": "'"$nome"'",
      "senha": "'"$senha"'"
    }'
}

# Função para criar conta CJ sequencialmente
create_account_cj() {
  cpf1="$1"
  cpf2="$2"
  senha="$3"
  curl -s -X POST \
    http://localhost:65501/criarContaCJ \
    -H 'Content-Type: application/json' \
    -d '{
      "cpf1": "'"$cpf1"'",
      "cpf2": "'"$cpf2"'",
      "senha": "'"$senha"'"
    }'
}

# Criação sequencial de contas PF
create_account_pf "CPF1" "Nome1" "senha123"
create_account_pf "CPF2" "Nome2" "senha456"
create_account_pf "CPF3" "Nome3" "senha789"

# Criação sequencial de contas PJ
create_account_pj "CNPJ1" "Empresa1" "senha123"
create_account_pj "CNPJ2" "Empresa2" "senha456"
create_account_pj "CNPJ3" "Empresa3" "senha789"

# Criação sequencial de contas CJ
create_account_cj "CPF1" "CPF2" "senha123"
create_account_cj "CPF2" "CPF3" "senha456"
create_account_cj "CPF3" "CPF1" "senha789"

# Verificação de contas duplicadas (apenas para garantia, no caso de algo inesperado)
# Não há necessidade de verificar duplicatas aqui, pois o teste é sequencial.

echo "Teste de criação de contas sequencial concluído."
