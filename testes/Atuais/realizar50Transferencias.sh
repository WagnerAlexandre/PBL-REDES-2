#!/bin/bash

# Base URL
BASE_URL="http://localhost:65502"

# Função para iniciar transferência
initiate_transfer() {
    local transactions="$1"
    response=$(curl -s -X POST "$BASE_URL/realizarTransferencia" -H "Content-Type: application/json" -d "$transactions")
    echo "Transfer Response: $response" >> transfer.log
}

# Limpar o arquivo de log
echo "" > transfer.log

# Criar um array de transações
create_transactions() {
    local count=$1
    local transactions="["

    for ((i = 1; i <= count; i++)); do
        local src1=$((i % 9 + 1))
        local src2=$(((i + 1) % 9 + 1))
        local dst=$(((i + 2) % 9 + 1))
        
        if [[ $i -ne 1 ]]; then
            transactions+=","
        fi
        
        transactions+=$(cat <<-END
        {
            "numConta": $src1,
            "valor": 10.0,
            "banco": "BG",
            "tipo": 2
        },
        {
            "numConta": $src2,
            "valor": 10.0,
            "banco": "BG",
            "tipo": 2
        },
        {
            "numConta": $dst,
            "valor": 20.0,
            "banco": "BG",
            "tipo": 1
        }
END
)
    done

    transactions+="]"
    echo "$transactions"
}

# Realizar 50 transferências entre contas
transactions=$(create_transactions 50)
initiate_transfer "$transactions"

echo "50 transferências realizadas entre as contas." >> transfer.log
