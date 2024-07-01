#!/bin/bash

# TESTES PARA CONCORRENCIA PARA DUAS CONTAS DIFERENTES: RESULTADO FINAL ESPERADO EM AMBAS: 0

# Definir as informações da conta para login
CPFRAZAO_PESSOA="11111111111"
SENHA_PESSOA="senha1"

CPFRAZAO_EMPRESA="11111111000100"
SENHA_EMPRESA="senha1"

# Função para realizar operações de depósito
depositar() {
    local cpfcnpj=$1
    local senha=$2
    local valor=$3
    local transacao="{\"cpfcnpj\": \"$cpfcnpj\", \"senha\": \"$senha\", \"valor\": $valor}"
    curl -X POST http://localhost:65501/somaLocal -H "Content-Type: application/json" -d "$transacao"
}

# Função para realizar operações de saque (valor negativo)
sacar() {
    local cpfcnpj=$1
    local senha=$2
    local valor=$3
    local valorNegativo=$(echo "$valor")  # Inverter o sinal para negativo usando bc
    local transacao="{\"cpfcnpj\": \"$cpfcnpj\", \"senha\": \"$senha\", \"valor\": $valorNegativo}"
    curl -X POST http://localhost:65501/somaLocal -H "Content-Type: application/json" -d "$transacao"
}

# Realizar 25 depósitos de 127 reais e 25 saques de 127 reais em paralelo
for i in {1..25}; do
    depositar $CPFRAZAO_PESSOA $SENHA_PESSOA 127 & sacar $CPFRAZAO_PESSOA $SENHA_PESSOA 127 &
    depositar $CPFRAZAO_EMPRESA $SENHA_EMPRESA 127 & sacar $CPFRAZAO_EMPRESA $SENHA_EMPRESA 127 &
done

wait  # Esperar todas as operações paralelas terminarem

echo "Operações concluídas."
read -p "Qualquer tecla para prosseguir"