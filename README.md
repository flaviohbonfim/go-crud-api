# Go CRUD API

Este é um projeto de API RESTful em Go, projetado para gerenciar usuários e produtos. Ele segue as melhores práticas de arquitetura, utilizando uma estrutura de projeto limpa e modular, com foco em separação de responsabilidades e testabilidade.

## Funcionalidades

- **Gerenciamento de Usuários:**
    - Criação, leitura, atualização e exclusão de usuários.
    - Autenticação de usuários via JWT.
- **Gerenciamento de Produtos:**
    - Criação, leitura, atualização e exclusão de produtos.
- **API RESTful:**
    - Endpoints bem definidos para cada recurso.
    - Respostas padronizadas em JSON.
- **Banco de Dados:**
    - Integração com PostgreSQL.
    - Migrações de banco de dados usando `migrate`.
- **Docker:**
    - Configuração para rodar a aplicação e o banco de dados em contêineres Docker.
- **Documentação:**
    - Documentação da API gerada automaticamente com Swagger.
- **Logging:**
    - Sistema de logging configurável.

## Tecnologias Utilizadas

- **Go:** Linguagem de programação principal.
- **PostgreSQL:** Banco de dados relacional.
- **Docker & Docker Compose:** Para orquestração de contêineres.
- **JWT (JSON Web Tokens):** Para autenticação.
- **Swagger:** Para documentação da API.
- **Viper:** Para gerenciamento de configurações.
- **GORM:** ORM para interação com o banco de dados (ou similar, dependendo da implementação).
- **BCrypt:** Para hashing de senhas.

## Estrutura do Projeto

```
.
├── cmd/                # Ponto de entrada da aplicação
│   └── api/            # Aplicação principal da API
│       └── main.go
├── docs/               # Documentação Swagger
├── internal/           # Código interno da aplicação
│   ├── config/         # Configurações da aplicação
│   ├── database/       # Conexão e migrações do banco de dados
│   ├── domain/         # Lógica de negócio (entidades, serviços, repositórios, handlers)
│   │   ├── products/
│   │   └── users/
│   ├── http/           # Camada HTTP (rotas, middlewares, handlers)
│   │   └── middleware/
│   ├── logger/         # Configuração de logging
│   └── repository/     # Implementações de repositórios
├── migrations/         # Arquivos de migração do banco de dados
├── pkg/                # Pacotes utilitários e genéricos
│   ├── jwt/            # Funções JWT
│   ├── pagination/     # Lógica de paginação
│   ├── password/       # Funções de hashing de senha
│   └── web/            # Utilitários web (respostas HTTP)
├── .env.example        # Exemplo de arquivo de variáveis de ambiente
├── docker-compose.yml  # Configuração Docker Compose
├── go.mod              # Módulos Go
├── go.sum              # Checksums de módulos Go
├── Makefile            # Comandos úteis para desenvolvimento
└── README.md           # Este arquivo
```

## Como Rodar o Projeto

### Pré-requisitos

- Go (versão 1.18 ou superior)
- Docker e Docker Compose
- Make (opcional, para usar os comandos do Makefile)

### Configuração

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/seu-usuario/go-crud-api.git
    cd go-crud-api
    ```

2.  **Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto, baseado no `.env.example`.
    ```bash
    cp .env.example .env
    ```
    Edite o arquivo `.env` com suas configurações, especialmente as credenciais do banco de dados e a chave JWT.

### Usando Docker Compose (Recomendado)

A maneira mais fácil de rodar o projeto é usando Docker Compose, que irá configurar o banco de dados e a aplicação Go.

1.  **Construa e inicie os contêineres:**
    ```bash
    docker-compose up --build
    ```
    Isso irá construir a imagem Go, iniciar o contêiner do PostgreSQL e o contêiner da API.

2.  **Executar Migrações (se necessário):**
    As migrações são aplicadas automaticamente na inicialização do contêiner da API. Se precisar rodar manualmente:
    ```bash
    docker-compose run --rm api migrate up
    ```

3.  **Acessar a API:**
    A API estará disponível em `http://localhost:8080`.
    A documentação Swagger estará disponível em `http://localhost:8080/swagger/index.html`.

### Rodando Localmente (Sem Docker para a Aplicação Go)

Se você preferir rodar a aplicação Go diretamente em sua máquina, mas ainda usar o Docker para o PostgreSQL:

1.  **Inicie apenas o contêiner do PostgreSQL:**
    ```bash
    docker-compose up -d postgres
    ```

2.  **Execute as migrações:**
    Certifique-se de ter a ferramenta `migrate` instalada (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`).
    ```bash
    migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
    ```
    (Ajuste a string de conexão conforme seu `.env`)

3.  **Inicie a aplicação Go:**
    ```bash
    go run cmd/api/main.go
    ```

## Endpoints da API

A documentação completa da API pode ser encontrada em `http://localhost:8080/swagger/index.html` após a aplicação estar rodando.

### Exemplos de Endpoints:

-   **Auth:**
    -   `POST /api/v1/auth/register` - Registrar um novo usuário
    -   `POST /api/v1/auth/login` - Autenticar um usuário
-   **Users:**
    -   `GET /api/v1/users` - Listar todos os usuários (requer autenticação)
    -   `GET /api/v1/users/{id}` - Obter um usuário por ID (requer autenticação)
    -   `PUT /api/v1/users/{id}` - Atualizar um usuário (requer autenticação)
    -   `DELETE /api/v1/users/{id}` - Excluir um usuário (requer autenticação)
-   **Products:**
    -   `GET /api/v1/products` - Listar todos os produtos
    -   `GET /api/v1/products/{id}` - Obter um produto por ID
    -   `POST /api/v1/products` - Criar um novo produto (requer autenticação)
    -   `PUT /api/v1/products/{id}` - Atualizar um produto (requer autenticação)
    -   `DELETE /api/v1/products/{id}` - Excluir um produto (requer autenticação)

## Testes

Para rodar os testes da aplicação:

```bash
go test ./...
```

## Contribuição

Sinta-se à vontade para contribuir com este projeto. Por favor, siga as diretrizes de contribuição (se houver) e abra um Pull Request.

## Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
