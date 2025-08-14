# Projeto Go — API CRUD (Usuários & Produtos) com Postgres, JWT, Testes e Boas Práticas

## Objetivo

Construir uma API REST em Go para cadastro e autenticação de usuários e CRUD de produtos, usando PostgreSQL como banco de dados e JWT para autenticação. O projeto deve ser simples para quem está começando em Go, mas aplicar um conjunto sólido de melhores práticas de arquitetura, testes, qualidade e automação.

## Stack e Principais Dependências

* **Go**: 1.22+
* **Framework HTTP**: `chi` (leve e idiomático)
* **ORM/DB**: `gorm.io/gorm` + `gorm.io/driver/postgres` (ou `sqlx` se preferir SQL puro; escolha uma)
* **Migrações**: `golang-migrate`
* **JWT**: `github.com/golang-jwt/jwt/v5`
* **Hash de senha**: `golang.org/x/crypto/bcrypt`
* **Validação**: `github.com/go-playground/validator/v10`
* **Config**: `github.com/spf13/viper`
* **Logs**: `github.com/rs/zerolog`
* **Docs (OpenAPI)**: `github.com/swaggo/swag` + `github.com/swaggo/http-swagger`
* **Lint**: `golangci-lint`
* **Testes**: `testing`, `httptest`, `testify` (assert/require)

## Arquitetura (Clean/Hexa simplificada)

Camadas e responsabilidades:

* **/cmd/api**: `main.go` inicializa config, logger, DB, rotas e servidor.
* **/internal/config**: carga de variáveis de ambiente com Viper.
* **/internal/logger**: configuração do zerolog.
* **/internal/database**: conexão e migrações.
* **/internal/domain/{users,products}**:

  * `entity.go` (modelos/entidades)
  * `repository.go` (interfaces)
  * `service.go` (regras de negócio)
  * `handler.go` (HTTP handlers)
* **/internal/repository**: implementações GORM/SQL das interfaces.
* **/internal/http**:

  * `router.go` (chi, middlewares)
  * `middleware/` (auth JWT, requestID, recovery, CORS)
  * `response.go` (respostas padronizadas, erros)
* **/pkg**: utilitários (jwt, password, pagination).
* **/migrations**: arquivos SQL de migração (up/down).
* **/docs**: arquivos gerados da OpenAPI (via swag).

**Padrões aplicados**

* **Repository Pattern** para desacoplar persistência.
* **Service Layer** para regras de negócio, transações e validações.
* **DTOs** (requests/responses) no handler para separar entidade de payload.
* **Erro estruturado** (códigos de erro + mensagens amigáveis) e mapeamento para HTTP.
* **Context** em todas as operações (cancelamento, deadlines).
* **Dependency Injection simples** via construtores (sem framework).

## Modelos (MVP)

### User

* `ID (uuid)`
* `Name (string, 2–100)`
* `Email (string, unique, válido)`
* `PasswordHash (string)`
* `Role (enum: "user"|"admin")`
* `CreatedAt/UpdatedAt (timestamp)`

### Product

* `ID (uuid)`
* `Name (string, 2–120)`
* `Description (string, opcional)`
* `Price (decimal >= 0)`
* `Stock (int >= 0)`
* `OwnerID (uuid, FK -> users.id)`
* `CreatedAt/UpdatedAt`

## Regras de Negócio

* **Auth**: login por email+senha → retorna `access_token` (JWT, exp. 15m) e `refresh_token` (exp. 7d).
* **Permissões**:

  * `admin`: CRUD de qualquer produto e listar usuários.
  * `user`: CRUD apenas dos **seus** produtos; não pode listar usuários.
* **Senhas**: sempre com `bcrypt` (cost 10–12).
* **Email**: único e case-insensitive.
* **Soft deletes**: não necessário no MVP (pode ser um TODO).

## Endpoints (REST)

Auth:

* `POST /v1/auth/register` → cria usuário (público)
* `POST /v1/auth/login` → retorna tokens (público)
* `POST /v1/auth/refresh` → novo access token (público)
* `POST /v1/auth/logout` → invalida refresh (opcional: blacklist/versão de token)

Users (admin):

* `GET /v1/users` → paginação, filtro por email
* `GET /v1/users/{id}`
* `DELETE /v1/users/{id}` (opcional)

Products:

* `GET /v1/products` → paginação (`page`, `page_size`), filtros (`name`, faixa de `price`), ordenação (`sort=name|price|-price`)
* `GET /v1/products/{id}`
* `POST /v1/products` (auth)
* `PUT /v1/products/{id}` (auth; owner/admin)
* `DELETE /v1/products/{id}` (auth; owner/admin)

## Autenticação & Segurança

* **JWT** assinado com HS256; `sub` = userID, `role` em `claims`.
* **Middleware**:

  * `RequestID`, `Recovery` (panic safe), `Logger`, `CORS`, `Auth` (checa header `Authorization: Bearer <token>`).
* **Headers seguros** (ex.: `X-Content-Type-Options`, `X-Frame-Options`, `Content-Security-Policy` — via middleware simples).
* **Rate limit** (TODO opcional: `golang.org/x/time/rate`).
* **Armazenar** apenas hash de senha; nunca retornar campos sensíveis.

## Validação & Respostas

* Validar payloads com `validator` + tags (`validate:"required,email"` etc.).
* Respostas JSON padronizadas:

  ```json
  { "data": {...}, "error": null, "meta": { "request_id": "..." } }
  { "data": null, "error": { "code": "bad_request", "message": "..." } }
  ```

## Banco de Dados & Migrações

* **Docker-compose** com Postgres 16.
* **Migrations** com `golang-migrate` (`migrations/0001_create_users.up.sql` etc.).
* Índices: `users(email unique)`, `products(name)`, `products(owner_id)`, `products(price)`.

## Configuração (12-factor)

* `.env` (exemplo):

  ```
  APP_NAME=go-crud
  APP_ENV=dev
  HTTP_PORT=8080
  JWT_SECRET=super-secret-change-me
  ACCESS_TOKEN_TTL=15m
  REFRESH_TOKEN_TTL=168h

  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=postgres
  DB_NAME=go_crud
  DB_SSLMODE=disable
  ```
* Viper carrega env + defaults; **nunca** fazer commit do `.env` real.

## Qualidade, Testes e CI

* **Testes unitários**:

  * Serviços (regras de negócio) com mocks do repositório.
  * Utilitários (hash, jwt, validação).
* **Testes de handler**:

  * `httptest` para rotas chave (login, criar produto, autorização).
* **Cobertura**: meta inicial 70%+.
* **golangci-lint**: vet, errcheck, staticcheck, revive, gofmt, gocyclo.
* **CI (GitHub Actions)**:

  * Jobs: `lint`, `test`, `build`.
  * Cache de módulos Go.
  * Rodar migrações e testes que dependem de Postgres com `services: postgres`.

## Observabilidade

* **Logs estruturados** (JSON) com `zerolog` (correlacionados por `request_id`).
* **Métricas** (TODO opcional): expor `/metrics` (Prometheus) com contadores de requisições e latência.
* **Healthcheck**: `GET /healthz` (checa DB com `ping`).

## Docker & Makefile

* **docker-compose.yml** (API + Postgres).
* **Makefile** com alvos:

  * `make dev` (rodar API com `air` ou `go run ./cmd/api`)
  * `make lint`, `make test`, `make cover`
  * `make migrate-up`, `make migrate-down`, `make migrate-create name=...`
  * `make build` (gera binário)
  * `make run` (usa binário)

## Documentação (OpenAPI)

* Anotar handlers com comentários do `swag` e gerar docs:

  * `swag init -g cmd/api/main.go -o ./docs`
* Expor `GET /swagger/*` com `http-swagger`.

## Fluxos Principais (Critérios de Aceite)

1. **Registro de usuário**

   * Input: `name`, `email`, `password` (>= 8 chars).
   * Output: usuário sem `passwordHash`.
   * DB: cria `role="user"` por padrão.
   * Valida e-mail único.

2. **Login**

   * Input: `email`, `password`.
   * Output: `access_token` (15m), `refresh_token` (7d).
   * Teste: senha incorreta → `401`.

3. **Criar produto (user)**

   * Requer `Authorization: Bearer <access_token>`.
   * Define `OwnerID` = `sub` do token.
   * Retorna `201` com produto criado.

4. **Atualizar produto (owner/admin)**

   * Owner pode atualizar seus produtos; `user` não owner → `403`.
   * `admin` pode atualizar qualquer produto.

5. **Listar produtos**

   * Paginado; filtros e ordenação.
   * Retorna `X-Total-Count` (ou em `meta`) para total.

6. **Listar usuários (admin)**

   * `user` → `403`; `admin` → `200` com paginação.

## Esqueleto de Pastas (exemplo)

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── logger/
│   ├── http/
│   │   ├── middleware/
│   │   ├── router.go
│   │   └── response.go
│   ├── domain/
│   │   ├── users/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   └── products/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── handler.go
│   └── repository/
│       ├── users_repo.go
│       └── products_repo.go
├── pkg/
│   ├── jwt/
│   ├── password/
│   └── pagination/
├── migrations/
├── docs/
├── .github/workflows/ci.yml
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Boas Práticas de Código

* Funções curtas, nomes claros, early-return para erros.
* Erros embrulhados com `fmt.Errorf("contexto: %w", err)`; nunca ignorar erro.
* Não exportar o que não precisa ser exportado.
* Handlers **somente** orquestram: validam input → chamam service → mapeiam resposta.
* Repositórios **somente** consultam DB; nenhum `json`/HTTP lá.
* Services contêm regras de negócio; transações (quando necessárias) gerenciadas aqui.
* Tests: um `Arrange-Act-Assert` limpo por caso; nomes de testes descritivos.

## Dicas de Implementação (passo a passo)

1. Inicie módulo: `go mod init github.com/<user>/go-crud-api`
2. Adicione dependências (chi, gorm, jwt, viper, validator, zerolog).
3. Crie `config` (envs) e `logger`.
4. Conecte Postgres e rode `migrate up`.
5. Modele entidades e repositórios.
6. Implemente serviços (users: register/login; products: CRUD).
7. Crie middlewares e rotas.
8. Escreva testes unitários de serviços e handlers críticos.
9. Gere OpenAPI e exponha `/swagger`.
10. Configure CI (lint + test).
11. Escreva README com **como rodar** (Docker e local).

## Exemplos de Critérios de Qualidade (QA)

* `golangci-lint` sem erros.
* `go test ./... -cover` ≥ 70%.
* `GET /healthz` retorna `200` e `db_status=ok`.
* Swagger acessível em `/swagger/index.html`.
* JWT inválido/expirado → `401`; sem escopo → `403`.

---