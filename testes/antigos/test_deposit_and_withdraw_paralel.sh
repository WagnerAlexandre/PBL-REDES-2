#!/bin/bash

# TESTE PARA TESTAR CONCORRENCIA EM UMA UNICA CONTA, RESULTADO FINAL ESPERADO: -500

# Definir as informações da conta para login
CPFRAZAO="11111111000100"
NUMCONTA=1
TIPO=2
SENHA="senha1"

# Função para realizar operações de depósito
depositar() {
    local valor=$1
    local transacao="{\"idConta\": $NUMCONTA, \"tipoconta\": $TIPO, \"valor\": $valor}"
    curl -X POST http://localhost:65501/somaLocal -H "Content-Type: application/json" -d "$transacao"
}

# Função para realizar operações de saque (valor negativo)
sacar() {
    local valor=$1
    local valorNegativo=$(echo "$valor")  # Inverter o sinal para negativo usando bc
    local transacao="{\"idConta\": $NUMCONTA, \"tipoconta\": $TIPO, \"valor\": $valorNegativo}"
    curl -X POST http://localhost:65501/somaLocal -H "Content-Type: application/json" -d "$transacao"
}

# Realizar 20 depósitos de 25 reais cada (em paralelo)
for i in {1..20}; do
    depositar 25 & sacar -50 &
done


wait  # Esperar todas as operações paralelas terminarem

echo "Saques concluídos."
read -p "Qualquer tecla para prosseguir"
