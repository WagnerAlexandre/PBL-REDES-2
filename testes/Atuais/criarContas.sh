#!/bin/bash

# Base URL
BASE_URL="http://localhost:65502"

# Função para criar uma conta PF
create_account_pf() {
    response=$(curl -s -X POST "$BASE_URL/criarContaPF" -H "Content-Type: application/json" -d '{
        "cpfcnpj": "'"$1"'",
        "nome": "'"$2"'",
        "senha": "password"
    }')
    echo "PF Account Created: $response" >> create_accounts.log
}

# Função para criar uma conta PJ
create_account_pj() {
    response=$(curl -s -X POST "$BASE_URL/criarContaPJ" -H "Content-Type: application/json" -d '{
        "cpfcnpj": "'"$1"'",
        "nome": "'"$2"'",
        "senha": "password"
    }')
    echo "PJ Account Created: $response" >> create_accounts.log
}

# Função para criar uma conta CJ
create_account_cj() {
    response=$(curl -s -X POST "$BASE_URL/criarContaCJ" -H "Content-Type: application/json" -d '{
        "cpf1": "'"$1"'",
        "cpf2": "'"$2"'",
        "senha": "password"
    }')
    echo "CJ Account Created: $response" >> create_accounts.log
}

# Limpar o arquivo de log
echo "" > create_accounts.log

# Criar 3 contas PF
create_account_pf "12345678901" "John Doe" &
create_account_pf "22345678901" "Jane Doe" &
create_account_pf "32345678901" "Jim Doe" &

# Criar 3 contas PJ
create_account_pj "12345678000100" "Company A" &
create_account_pj "22345678000100" "Company B" &
create_account_pj "32345678000100" "Company C" &

# Criar 3 contas CJ
create_account_cj "12345678901" "09876543210" &
create_account_cj "22345678901" "19876543210" &
create_account_cj "32345678901" "29876543210" &

wait
echo "Todas as contas foram criadas." >> create_accounts.log

# Exibir o log no terminal
cat create_accounts.log
