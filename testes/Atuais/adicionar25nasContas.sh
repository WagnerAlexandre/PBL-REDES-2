#!/bin/bash

# Base URL
BASE_URL="http://localhost:65502"

# Função para somar saldo
sum_balance() {
    response=$(curl -s -X POST "$BASE_URL/somaLocal" -H "Content-Type: application/json" -d '{
        "idConta": '"$1"',
        "tipoconta": '"$2"',
        "valor": 25.0
    }')
    echo "Sum Balance (Conta $1 Tipo $2): $response" >> add_balance.log
}

# Limpar o arquivo de log
echo "" > add_balance.log

# Adicionar R$ 25,00 em cada conta PF
sum_balance 1 1
sum_balance 2 1
sum_balance 3 1

# Adicionar R$ 25,00 em cada conta PJ
sum_balance 4 2
sum_balance 5 2
sum_balance 6 2

# Adicionar R$ 25,00 em cada conta CJ
sum_balance 7 3
sum_balance 8 3
sum_balance 9 3

echo "Saldo adicionado a todas as contas." >> add_balance.log
