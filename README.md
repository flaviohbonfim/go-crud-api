# Projeto Go — API CRUD (Usuários & Produtos)

## Visão Geral

Este projeto é uma API RESTful desenvolvida em Go, focada no gerenciamento de usuários e produtos. Ele incorpora as melhores práticas de arquitetura, testes, segurança e automação, sendo um excelente ponto de partida para quem deseja aprender Go construindo uma aplicação robusta.

## Objetivo

Construir uma API REST em Go para cadastro e autenticação de usuários e CRUD de produtos, usando PostgreSQL como banco de dados e JWT para autenticação. O projeto é projetado para ser didático, mas aplicando um conjunto sólido de melhores práticas de arquitetura, testes, qualidade e automação.

## Stack e Principais Dependências

*   **Go**: 1.22+
*   **Framework HTTP**: `chi` (leve e idiomático)
*   **ORM/DB**: `gorm.io/gorm` + `gorm.io/driver/postgres`
*   **Migrações**: `golang-migrate`
*   **JWT**: `github.com/golang-jwt/jwt/v5`
*   **Hash de senha**: `golang.org/x/crypto/bcrypt`
*   **Validação**: `github.com/go-playground/validator/v10`
*   **Config**: `github.com/spf13/viper`
*   **Logs**: `github.com/rs/zerolog`
*   **Docs (OpenAPI)**: `github.com/swaggo/swag` + `github.com/swaggo/http-swagger`
*   **Testes**: `testing`, `httptest`, `testify` (assert/require)

## Arquitetura (Clean/Hexagonal simplificada)

A arquitetura do projeto segue princípios de Clean Architecture/Arquitetura Hexagonal, dividindo as responsabilidades em camadas claras:

*   **/cmd/api**: Ponto de entrada da aplicação (`main.go`), responsável por inicializar configurações, logger, conexão com o DB, rotas e o servidor HTTP.
*   **/internal/config**: Gerencia o carregamento de variáveis de ambiente usando Viper.
*   **/internal/logger**: Configuração centralizada do Zerolog para logs estruturados.
*   **/internal/database**: Lida com a conexão ao PostgreSQL e a aplicação de migrações.
*   **/internal/domain/{users,products}**: Contém a lógica de negócio específica para usuários e produtos.
    *   `entity.go`: Define os modelos/entidades de dados.
    *   `repository.go`: Define as interfaces para operações de persistência.
    *   `service.go`: Implementa as regras de negócio, validações e orquestra as operações de repositório.
    *   `handler.go`: Contém os handlers HTTP que processam as requisições e chamam os serviços.
*   **/internal/repository**: Implementações concretas das interfaces de repositório, utilizando GORM para interagir com o banco de dados.
*   **/internal/http**:
    *   `router.go`: Configura o roteador Chi, define as rotas e aplica middlewares.
    *   `middleware/`: Contém middlewares HTTP (autenticação JWT, recuperação de panics, etc.).
    *   `health.go`: Handler para o endpoint de health check.
*   **/pkg**: Pacotes de utilitários genéricos e reutilizáveis (JWT, hashing de senha, respostas HTTP padronizadas).
*   **/migrations**: Arquivos SQL para as migrações do banco de dados (`.up.sql`).
*   **/docs**: Arquivos gerados da especificação OpenAPI (via `swag`).
*   **Padrões Aplicados**: Repository Pattern, Service Layer, DTOs (Requests/Responses), Erro Estruturado, Context, Dependency Injection simples.

## Modelos (MVP)

### User

*   `ID (uuid)`
*   `Name (string, 2–100)`
*   `Email (string, unique, válido)`
*   `PasswordHash (string)`
*   `Role (enum: "user"|"admin")`
*   `CreatedAt/UpdatedAt (timestamp)`

### Product

*   `ID (uuid)`
*   `Name (string, 2–120)`
*   `Description (string, opcional)`
*   `Price (decimal >= 0)`
*   `Stock (int >= 0)`
*   `OwnerID (uuid, FK -> users.id)`
*   `CreatedAt/UpdatedAt`

## Regras de Negócio

*   **Autenticação**: Login por email+senha retorna `access_token` (JWT, exp. 15m) e `refresh_token` (exp. 7d).
*   **Permissões**:
    *   `admin`: CRUD de qualquer produto e listar usuários.
    *   `user`: CRUD apenas dos **seus** produtos; não pode listar usuários.
*   **Senhas**: Sempre com `bcrypt` (cost 10–12).
*   **Email**: Único e case-insensitive.

## Endpoints (REST)

### Autenticação
*   `POST /v1/auth/register` → Cria usuário (público)
*   `POST /v1/auth/login` → Retorna tokens (público)

### Usuários (Admin)
*   `GET /v1/users` → Lista usuários (requer `admin` role)

### Produtos
*   `POST /v1/products` → Cria produto (requer autenticação)
*   `GET /v1/products/{id}` → Busca produto por ID (requer autenticação)
*   `PUT /v1/products/{id}` → Atualiza produto (requer autenticação, owner ou admin)
*   `DELETE /v1/products/{id}` → Deleta produto (requer autenticação, owner ou admin)

### Outros
*   `GET /healthz` → Verifica a saúde da aplicação e conexão com o DB.
*   `GET /swagger/*` → Interface da documentação OpenAPI (Swagger UI).

## Autenticação & Segurança

*   **JWT** assinado com HS256; `sub` = userID, `role` em `claims`.
*   **Middlewares**: `RequestID`, `RealIP`, `Recoverer`, `AuthMiddleware` (checa `Authorization: Bearer <token>`), `HasRoleMiddleware` (para controle de acesso baseado em role).
*   **Armazenar** apenas hash de senha; nunca retornar campos sensíveis.

## Validação & Respostas

*   Validação de payloads com `validator` e tags (`validate:"required,email"`).
*   Respostas JSON padronizadas via `pkg/web`:
    ```json
    { "data": {...}, "error": null }
    { "data": null, "error": { "code": "bad_request", "message": "..." } }
    ```

## Banco de Dados & Migrações

*   **PostgreSQL 16** via `docker-compose`.
*   **Migrations** com `golang-migrate`, aplicadas automaticamente na inicialização da aplicação em ambiente de desenvolvimento.
*   **Índices**: `users(email unique)`, `products(name)`, `products(owner_id)`, `products(price)`.

## Configuração (12-factor)

*   Variáveis de ambiente carregadas via `.env` (exemplo em `.env.example`).
*   `viper` para carregar as configurações.

## Qualidade, Testes e CI

*   **Testes unitários**: Serviços (regras de negócio) com mocks do repositório, utilitários (hash, jwt, validação).
*   **Testes de handler**: `httptest` para rotas chave (login, criar produto, autorização).
*   **golangci-lint**: Ferramenta de linting para garantir a qualidade do código.
*   **CI (GitHub Actions)**: Configuração para `lint`, `test` e `build` (ainda a ser implementada).

## Observabilidade

*   **Logs estruturados** (JSON) com `zerolog`.
*   **Healthcheck**: `GET /healthz` (checa DB com `ping`).

## Docker & Makefile

*   **`docker-compose.yml`**: Define os serviços da API e do PostgreSQL.
*   **`Makefile`**: Contém alvos para facilitar o desenvolvimento e a automação.

## Documentação (OpenAPI)

*   Handlers anotados com comentários `swag`.
*   Documentação gerada via `swag init` e servida em `GET /swagger/*` via `http-swagger`.

## Como Rodar o Projeto

### Pré-requisitos

*   [Go](https://golang.org/dl/) (versão 1.22+)
*   [Docker](https://www.docker.com/get-started/) e [Docker Compose](https://docs.docker.com/compose/install/)
*   Um editor de código (VS Code, GoLand, etc.)

### Passos para Rodar

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/seu-usuario/go-crud-api.git
    cd go-crud-api
    ```

2.  **Configure as variáveis de ambiente:**
    Crie um arquivo `.env` na raiz do projeto, copiando o conteúdo de `.env.example` e ajustando conforme necessário.
    ```bash
    cp .env.example .env
    # Edite o .env se precisar mudar portas ou credenciais
    ```

3.  **Inicie o banco de dados:**
    ```bash
    docker-compose up -d postgres
    ```
    Aguarde alguns segundos para o PostgreSQL iniciar completamente.

4.  **Instale as dependências Go:**
    ```bash
    go mod tidy
    ```

5.  **Gere a documentação OpenAPI:**
    ```bash
    go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go -o ./docs
    ```

6.  **Inicie a aplicação:**
    ```bash
    go run ./cmd/api/main.go
    ```
    As migrações do banco de dados serão aplicadas automaticamente na inicialização.

7.  **Acesse a API:**
    *   **Swagger UI**: `http://localhost:8080/swagger/`
    *   **Health Check**: `http://localhost:8080/healthz`
    *   **Endpoints da API**: Use ferramentas como `curl` ou Postman para interagir com os endpoints de autenticação, usuários e produtos.

## Fluxos Principais (Critérios de Aceite)

*   **Registro de usuário**: `POST /v1/auth/register`
*   **Login**: `POST /v1/auth/login` (retorna `access_token` e `refresh_token`)
*   **Criação de produto**: `POST /v1/products` (requer `access_token`)
*   **Atualização de produto**: `PUT /v1/products/{id}` (requer `access_token`, owner ou admin)
*   **Listagem de usuários**: `GET /v1/users` (requer `access_token` de `admin`)

## Boas Práticas de Código

*   Funções curtas, nomes claros, early-return para erros.
*   Erros embrulhados com `fmt.Errorf("contexto: %w", err)`; nunca ignorar erro.
*   Não exportar o que não precisa ser exportado.
*   Handlers **somente** orquestram: validam input → chamam service → mapeiam resposta.
*   Repositórios **somente** consultam DB; nenhum `json`/HTTP lá.
*   Services contêm regras de negócio; transações (quando necessárias) gerenciadas aqui.
*   Tests: um `Arrange-Act-Assert` limpo por caso; nomes de testes descritivos.