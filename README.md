# Client-Server-API

Desafio realizado para o curso [Go-Expert](https://goexpert.fullcycle.com.br/pos-goexpert/)

## Descrição

O Projeto está separado em dois apps:

- **server**: recupera as informações da cotação do Dolar, salva os dados recuperados no Banco de dados SQLite, e retorna elas no formato JSON
-  **client**: recupera as informações do app server e salva a contação atual em um arquivo de nome predefinido

## Executar

### Server

Para executar o app server é necessário entrar na pasta `server` e rodar o seguinte comando:

```
go run main.go
```

### Client
Para executar o app client é necessário executar o app server, após isso entrar na pasta `client` e rodar o seguinte comando:

```
go run main.go
```