# PBL-REDES-2
## Sumário
 1. [Introdução](#introdução)
 2. [Como Utilizar o Sistema](#como-utilizar-o-sistema)
 3. [Estruturas Utilizadas Pelo Sistema](#estruturas-utilizadas-no-sistema)

_____
## Introdução
Este sistema bancário descentralizado foi minha solução proposta para o segundo problema da disciplina PLB de Concorrência e Conectividade - TEC502 da Universidade Estadual de Feira de Santana (UEFS) do semestre 24.1.

O problema pedia o desenvolvimento de um sistema descentralizado no qual um cliente pode utilizar o saldo de qualquer conta que o pertença (que seu CPF ou CNPJ esteja registrado como titular daquela conta) sendo esta do banco que ele está atualmente logado ou de qualquer outro banco participante do consorcio bancário para realizar transferências entre as contas dos bancos, e, estas transferências devem seguir um modelo de transação atômica para assegurar que não haja movimentação indevida dos saldos ou algum tipo de instabilidade que faça que o sistema realize uma transação indevida ou deixe uma operação incompleta.

Foi optado por desenvolver um consorcio de somente 3 bancos:
O Banco Bola Monetária Nacional (ou BBMN), o Banco Brasileirinho (ou BB) e o Banco Gringuesco (ou BG).
O modelo de transação atômica utilizado foi o 2-Phase-Commit (ou 2PC, para abreviar), o 2PC é baseado em dividir a transação em duas fases:
1. Fase de preparação: A fase de preparação envia um "pedido" que avisa os participantes para se prepararem para receber uma requisição, caso todos os participates avisem que estão prontos para receber a requisição, o algoritmo segue para a fase de commit. Caso algum participante retorne algum erro, o algoritmo entra em uma fase intermediaria, chamada de "abort".
2. Fase de Commit: Quando todos os participantes estão prontos para receber a requisição, o "atuador" (servidor/cliente) de origem envia a requisição (ou commit), e os participantes guardam ou atualizam os respectivos valores condizentes com o commit recebido. 
3. Fase de Abort: Quando algum participante retorna um erro, o algoritmo precisa cancelar todas as operações, esta fase envia um pedido de abort para todos os participantes. Esta fase garante atomicidade ao buscar garantir que sempre todos os participantes devem retornar um status de "Preparado".
----------------------------
## Como utilizar o sistema: 
### Requisitos:
O sistema foi desenvolvido utilizando Golang 1.22.2, então possuir o Golang no mínimo nesta versão
### Executando
Construa e suba as imagens docker para os containers e execute cada sistema bancário em um container.
Para construir a imagem docker, vá até a pasta que guarda um sistema bancário e digite no terminal:

` docker build . -t nome-do-banco `

Os bancos utilizam as portas para se diferenciarem, caso execute em computadores diferentes, será necessário alterar algumas FLAGS e linhas do código para que os códigos consigam acessar as rotas devidamente.

Para o BBMN:

`docker run --name nome-do-container --rm -p 65500:65500 nome-da-imagem`

Para o BB:
`docker run --name nome-do-container --rm -p 65502:65502 nome-da-imagem `

Para o BG:
`docker run --name nome-do-container --rm -p 65501:65501 nome-da-imagem `

### Acesso ao banco
No navegador, digite o ip:porta do banco desejado seguido de '/inicial' para ter acesso a tela de login.
Exemplo:
`https://localhost:65502/inicial`

Abaixo um fluxograma das paginas.
<img src="src/readme/loginFlux.png" alt="fluxoDoSistema" >

---------------------
## Estruturas utilizadas no sistema

### Estruturas utilizadas

Diagrama Geral de classes utilizadas no sistema.
<img src="src/readme/classes.png" alt="diagramaClasses" >

#### Classes de uso via Web
Estas classes são utilizadas para receber JSON'S pelas rotas acessadas pelas paginas web.

##### Login
<p align="center">
<img src="src/readme/web/classLogin.png" alt="classeLogin" >
<figurecaption #align="center"> Classe de Login</figurecaption>
</p>

Esta estrutura é responsável por receber o json da requisição de login.

- CPFRAZAO - Guarda o CPF ou CNPJ que está a logar
- NumConta - O Numero (ID) daquela conta.
- TIPO - Guarda o tipo daquela conta, seja 1 uma conta de pessoa física e 2 Pessoa Jurídica
- Senha - A senha utilizada para login.
Classes e seus usos

##### Conta
Resposta da requisição "/getContas" que procura todas as contas associadas a um CPF ou CNPJ.
<p align="center">
<img src="src/readme/web/classConta.png" alt="classeLogin" >
<figurecaption #align="center"> Classe de resposta ao solicitar dados da conta</figurecaption>
</p>

- ID - Representa o numero da conta.
- Nome - Nome Titular ou Empresarial da conta.
- Balanco - Saldo disponível na conta.
- Banco - Banco o qual esta conta pertence.

##### Transação
Utilizada internamente para realizar uma modificação local do saldo de uma conta.
<p align="center">
<img src="src/readme/web/classTransacao.png" alt="transacaoLocal" >
<figurecaption #align="center"> Classe de transação local</figurecaption>
</p>

- ID - Numero da conta a receber a transação.
- TipoConta - Tipo da conta a receber a transação (PF:1, PJ:2, CJ:3)
- Valor - Quantia a ser reduzida/incrementada.


##### Transação Web
Utilizada para receber um pedido de transação via rota.
<p align="center">
<img src="src/readme/web/classTransacaoWeb.png" alt="transacaoWeb" >
<figurecaption #align="center"> Classe de transação web</figurecaption>
</p>

- ID - Numero da conta alvo da transação.
- Valor - Valor a ser incrementado/decrementado do saldo.
- Banco - Banco o qual a conta pertence
- Tipo - Tipo da transação (1: Somar, 2: Decrementar)

#### Estruturas de uso interno
Estas estruturas são utilizadas somente internamente pelo servidor para representar os tipos de conta que é possível possuir no sistema.

##### Conta Pessoa Física
Como as contas de Pessoa Física são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaPF.png" alt="classPF" >
<figurecaption #align="center"> Classe de conta Pessoa Física</figurecaption>
</p>

- NumConta - Representa o número identificador da conta
- CPF   - Representa o CPF titular da conta.
- Nome - Nome do titular da conta.
- Senha - Senha da conta.
- Tipo - Tipo da conta, utilizado para passar para a struct Conta, ao enviar as informações da conta.
- Balanco - Saldo disponível na conta.
- mutex - Objeto de sincronização para tratar concorrência, ao utilizar mutex.lock(), a conta fica bloqueada, permitindo somente que a função que bloqueou a conta modifique o saldo.


##### Conta Pessoa Jurídica
Como as contas de Pessoa Jurídica são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaPJ.png" alt="classPJ" >
<figurecaption #align="center"> Classe de conta Pessoa Jurídica</figurecaption>
</p>

- NumConta - Representa o número identificador da conta
- CNPJ   - Representa o CNPJ titular da conta.
- Nome - Nome do Empresa titular da conta.
- Senha - Senha da conta.
- Tipo - Tipo da conta, utilizado para passar para a struct Conta, ao enviar as informações da conta.
- Balanco - Saldo disponível na conta.
- mutex - Objeto de sincronização para tratar concorrência, ao utilizar mutex.lock(), a conta fica bloqueada, permitindo somente que a função que bloqueou a conta modifique o saldo.

##### Conta Conjunta
Como as contas conjuntas são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaCj.png" alt="classCJ" >
<figurecaption #align="center"> Classe de conta conjunta</figurecaption>
</p>

- NumConta - Representa o número identificador da conta
- CPF1   - Representa o CPF  de um titular da conta.
- CPF2   - Representa o CPF  do outro titular da conta.
- Nome - Nome titular da conta.
- Senha - Senha da conta.
- Tipo - Tipo da conta, utilizado para passar para a struct Conta, ao enviar as informações da conta.
- Balanco - Saldo disponível na conta.
- mutex - Objeto de sincronização para tratar concorrência, ao utilizar mutex.lock(), a conta fica bloqueada, permitindo somente que a função que bloqueou a conta modifique o saldo.

### Estruturas utilizadas pelo 2-Phase-Commit (2PC) implementado
Estas estruturas são implementadas e utilizadas somente pelas rotas responsáveis por gerencias o 2PC

##### Prepare Request
Refere a primeira fase do 2PC, a fase de preparação, esta estrutura é enviada via a rota "/prepare" para os bancos participantes de uma transação para preparação das contas participantes.
<p align="center">
<img src="src/readme/2PC/classPrepareRequest.png" alt="classPrepararReq" >
<figurecaption #align="center"> Classe de preparar requisição</figurecaption>
</p>



##### Prepare Response
Esta estrutura é a resposta a requisição da fase de preparação do 2PC
<p align="center">
<img src="src/readme/2PC/classPrepareResponse.png" alt="classRespReq" >
<figurecaption #align="center"> Classe de resposta da fase de preparação</figurecaption>
</p>

### Implementação do modelo de transação atômica
O modelo implementado é o 2-Phase-Commit

<p align="center">
<img src="src/readme/2PC/sequencia2pc.png" alt="classRespReq" >
<figurecaption #align="center"> Diagrama de sequência do 2PC implementado</figurecaption>
</p>


-------------------------
## Conclusão
### 1. Permite gerenciar contas?
Sim, o sistema realiza o gerenciamento de contas, permitindo a criação de contas dos tipos Pessoa Física (PF), Pessoa Jurídica (PJ) e Conta Conjunta (CJ). As informações de cada conta, como número da conta, CPF/CNPJ, nome, senha e balanço, são armazenadas e gerenciadas no sistema.

### 2. Permite selecionar e realizar transferência entre diferentes contas?
Sim, o sistema permite selecionar e realizar transferências entre diferentes contas. As transferências podem ser realizadas tanto entre contas do mesmo banco quanto entre contas de diferentes bancos.

### 3. É possível transacionar entre diferentes bancos?
Sim, é possível transacionar entre diferentes bancos. O sistema permite enviar transações do banco A, B e C para o banco D, utilizando comunicação adequada entre os servidores dos bancos envolvidos.

### 4. Comunicação entre servidores
Os bancos estão se comunicando utilizando o protocolo HTTP, com endpoints específicos para preparar, realizar commit's e abortar transações, implementando o algoritmo de Commit em Duas Fases (2PC).

### 5. Sincronização em um único servidor
A concorrência em um único servidor é tratada utilizando mutexes (locks) para garantir que apenas uma transação pode alterar o saldo de uma conta por vez. Isso evita condições de corrida e garante a consistência dos dados.

### 6. Algoritmo da concorrência distribuída está teoricamente bem empregado? Qual algoritmo foi utilizado? Está correto para a solução?
Sim, o algoritmo de Commit em Duas Fases (2PC) foi utilizado para garantir a consistência e atomicidade das transações distribuídas. Este algoritmo é adequado para a solução, pois permite que todos os bancos participantes concordem em commitar ou abortar uma transação de forma coordenada.

### 7. Algoritmo está tratando o problema na prática? A implementação do algoritmo está funcionando corretamente?
Sim, a implementação do algoritmo 2PC está funcionando corretamente. O sistema realiza a preparação das transações em todos os bancos participantes e somente realiza o commit se todos os participantes confirmarem a preparação com sucesso. Caso contrário, a transação é abortada em todos os participantes.

### 8. Tratamento da confiabilidade
O sistema é projetado para continuar funcionando corretamente quando um dos bancos perde a conexão. Se um banco participante não puder ser alcançado, a transação é abortada. Quando o banco retorna à conexão, novas transações podem ser iniciadas.

### 9. Pelo menos uma transação concorrente é realizada?
Sim, o sistema suporta a realização de transações concorrentes. Utilizando mutexes, garantimos que múltiplas transações podem ser processadas de forma concorrente, mantendo a consistência dos saldos das contas.

### 10. Como foi tratado o caso em que mais de duas transações ocorrem no mesmo banco de forma concorrente? O saldo fica correto? Os clientes conseguem realizar as transações?
Quando mais de duas transações ocorrem no mesmo banco de forma concorrente, mutexes são utilizados para garantir que apenas uma transação pode alterar o saldo de uma conta por vez. Isso assegura que o saldo das contas permanece correto e que todas as transações são processadas com sucesso, mantendo a consistência dos dados.
