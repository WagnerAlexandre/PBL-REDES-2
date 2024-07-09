#!/bin/bash

# Base URL
BASE_URL="http://localhost:65502"

# Função para reduzir saldo
reduce_balance() {
    response=$(curl -s -X POST "$BASE_URL/reducaoLocal" -H "Content-Type: application/json" -d '{
        "idConta": '"$1"',
        "tipoconta": '"$2"',
        "valor": 15.0
    }')
    echo "Reduce Balance (Conta $1 Tipo $2): $response" >> reduce_balance.log
}

# Limpar o arquivo de log
echo "" > reduce_balance.log

# Reduzir R$ 15,00 em cada conta PF
reduce_balance 1 1
reduce_balance 2 1
reduce_balance 3 1

# Reduzir R$ 15,00 em cada conta PJ
reduce_balance 4 2
reduce_balance 5 2
reduce_balance 6 2

# Reduzir R$ 15,00 em cada conta CJ
reduce_balance 7 3
reduce_balance 8 3
reduce_balance 9 3

echo "Saldo reduzido de todas as contas." >> reduce_balance.log
