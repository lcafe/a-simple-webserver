# A simple webserver

Uma simples implementação de um webserver em Go.

## Configuração

A configuração é feita através do arquivo `internal/config/config.json`.

EXEMPLO: 
```json
{
    "http_port": ":8000",
    "proxy_prefix": "/app/",
    "virtual_hosts": {
        "/example" : "http://example.com",
        "/hello" : "http://example.org"
    },
    "file_server_root_url" : "/files/",
    "file_server":  "./public-example"
}
```

## Comandos

Para rodar o servidor, execute o comando `go run cmd/server/main.go`.

## Rotas

O servidor suporta as seguintes rotas:

- `/`: retorna "Hello World"
- `proxy_prefix`: redireciona para o VirtualHost definido no arquivo de configuração
- `file_server_root_url`: redireciona para o diretório definido no arquivo de configuração

## VirtualHosts

Para cada VirtualHost definido no arquivo de configuração, o servidor cria um proxy reverso.

Exemplo:

- Se o VirtualHost for `http://example.com`, o servidor criará um proxy reverso para `http://your-server:HTTP_PORT/proxy_prefix/example`.

## Arquivos estáticos

Os arquivos estáticos são servidos pelo servidor, sem a necessidade de redirecionamento.

## Escopo do projeto

O projeto está dividido em duas pastas:

- `cmd/server`: comandos para rodar o servidor
- `internal`: código do servidor

### cmd/server

O comando `go run cmd/server/main.go` executa o servidor.

### internal

- `config`: arquivo de configuração
- `handlers`: implementação dos handlers

#### config

Arquivo de configuração do servidor.

- `config.go`: implementa a lógica de carregamento da configuração
- `config.json`: arquivo de configuração do servidor

#### handlers

Implementação dos handlers do servidor. 

- `file.go`: implementa o handler de arquivos estáticos
- `proxy.go`: implementa o handler de VirtualHosts

### A intenção do projeto

Este projeto foi criado para fins de aprendizado e testes. Não é recomendado para uso em produção.

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.