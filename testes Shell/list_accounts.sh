#!/bin/bash

# URL base do servidor
BASE_URL="http://localhost:65501"

# Função para listar contas PF
list_accounts_pf() {
  echo "Listando contas PF..."
  curl -s -X GET "$BASE_URL/contasPF" 
}

# Função para listar contas PJ
list_accounts_pj() {
  echo "Listando contas PJ..."
  curl -s -X GET "$BASE_URL/contasPJ" 
}

# Função para listar contas CJ
list_accounts_cj() {
  echo "Listando contas CJ..."
  curl -s -X GET "$BASE_URL/contasCJ" 
}

# Executa as funções de listagem
list_accounts_pf
list_accounts_pj
list_accounts_cj

echo "Listagem de contas finalizada."
read -p "Pressione Enter para fechar..."