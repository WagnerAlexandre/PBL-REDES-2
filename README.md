# PBL-REDES-2
## Sumário
 1. [Introdução](#introdução)
 2. [Como Utilizar o Sistema](#como-utilizar-o-sistema)
 3. [Estruturas Utilzadas Pelo Sistema](#estruturas-utilizadas-no-sistema)

_____
## Introdução
Este sistema bancário descentralizado foi minha solução proposta para o segundo problema da disciplina PLB de Concorrência e Conectividade - TEC502 da Universidade Estadual de Feira de Santana (UEFS) do semestre 24.1.

O problema pedia o desenvolvimento de um sistema descentralizado no qual um cliente pode utilizar o saldo de qualquer conta que o pertença (que seu cpf ou cpnj esteja registrado como titular daquela conta) sendo esta do banco que ele está atualmente logado ou de qualquer outro banco participante do consorcio bancário para realizar transferências entre as contas dos bancos, e, estas transferências devem seguir um modelo de transação atômica para assegurar que não haja movimentação indevida dos saldos ou algum tipo de instabilidade que faça que o sistema realize uma transação indevida ou deixe uma operação incompleta.

Foi optado por desenvolver um consorcio de somente 3 bancos:
O Banco Bola Monetária Nacional (ou BBMN), o Banco Brasileirinho (ou BB) e o Banco Gringuesco (ou BG).

----------------------------
## Como utilizar o sistema: 
### Requisitos:
O sistema foi desenvolvido utilizando Golang 1.22.2, então possuir o Golang no mínimo nesta versão
### Executando
Construa e suba as imagens docker para os containers e execute cada sistema bancário em um container.
Para construir a imagem docker, vá até a pasta que guarda um sistema bancário e digite no terminal:

` docker build . -t nome-do-banco `

Os bancos utilizam as portas para se diferenciarem, caso execute em computadores diferentes, será necessário alterar algumas FLAGS e linhas do código para que os códigos consigam acessar as rotas devidamente.

Para o BNMN:

`docker run -p 65500:65500 --rm -it -e wagnerAlexandre/banco-bbmn
docker build -t banco-bbmn .`

Para o BG:
`docker run -p 65501:65501 --rm -it -e wagnerAlexandre/banco-bg
docker build -t banco-bg .`

Para o BB:
`docker run -p 65502:65502 --rm -it -e wagnerAlexandre/banco-bb
docker build -t banco-bb .`

### Acesso ao banco
No navegador, digite o ip:porta do banco desejado seguido de '/inicial' para ter acesso a tela de login.
Exemplo:
`https://0.0.0.0:65501/inicial`

Abaixo um fluxograma das paginas.
<img src="src/readme/loginFlux.png" alt="fluxoDoSistema" >

---------------------
## Estruturas utilizadas no sistema

### Estruturas utilizadas

Diagrama Geral de classes utilizadas no sistema.
<img src="src/readme/classes.png" alt="diagramaClasses" >

#### Classes de uso via Web
Estas classes são utilizadas para receber JSONS pelas rotas acessadas pelas paginas web.

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

##### Transacao
Utilizada internamente para realizar uma modificação local do saldo de uma conta
<p align="center">
<img src="src/readme/web/classTransacao.png" alt="transacaoLocal" >
<figurecaption #align="center"> Classe de transação local</figurecaption>
</p>

##### Transação Web
Utilizada para receber um pedido de transação via rota.
<p align="center">
<img src="src/readme/web/classTransacaoWeb.png" alt="transacaoWeb" >
<figurecaption #align="center"> Classe de transação web</figurecaption>
</p>

#### Estruturas de uso interno
Estas estruturas são utilizadas somente internamente pelo servidor.

##### Conta Pessoa Física
Como as contas de Pessoa Física são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaPF.png" alt="classPF" >
<figurecaption #align="center"> Classe de conta Pessoa Física</figurecaption>
</p>

##### Conta Pessoa Juridica
Como as contas de Pessoa Jurídica são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaPJ.png" alt="classPJ" >
<figurecaption #align="center"> Classe de conta Pessoa Jurídica</figurecaption>
</p>

##### Conta Conjunta
Como as contas conjuntas são registradas no sistema.
<p align="center">
<img src="src/readme/Banco de dados/classContaCj.png" alt="classCJ" >
<figurecaption #align="center"> Classe de conta conjunta</figurecaption>
</p>

### Estruturas utilizadas pelo 2-Phase-Commit (2PC) implementado
Estas estruturas são implementadas e utilizadas somente pelas rotas responsaveis por gerencias o 2PC

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
Os bancos estão se comunicando utilizando o protocolo HTTP, com endpoints específicos para preparar, commitar e abortar transações, implementando o algoritmo de Commit em Duas Fases (2PC).

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
