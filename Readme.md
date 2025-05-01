This it template that can be used to setup enterprise level application in Go

To generate docs use:
```shell
swag init
```

Starting Postgres:
```shell
docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:17-alpine
```

# Build go program

```shell
go build -trimpath .
```

# Required items:
1. [X] Create endpoint to handle simple CRUD requests
   - [X] Choose the right framework to implement (Used `gin`, although `chia` also can be used)
   - [X] Define project structure
   - [X] Implement functionality
   - [X] Add swagger (accessible by path localhost:8080/swagger/index.html )
2. [X] Create integration with Postgres database
   - [X] Choose the right library (`pgx` with `squirrel` can be used)
   - [X] Implement functionality
   - [X] Add database migration
3. [X] Setup logging
   - [X] Use slog to implement it (we need to supply module name, timestamp, severity, message)
4. [ ] Test it automatically
   - [ ] Add unit tests
   - [ ] Add integration tests
   - [ ] Add test coverage report
5. [ ] Healthcheck route
6. [ ] Infrastructure setup
   - [ ] Put application in Docker
   - [ ] Create `docker-compose.yaml` file for deployment
   - [ ] Create CI/CD pipeline and push to GitHub repository
   - [ ] Add code linters
   - [ ] Add logger linters ([vet](https://pkg.go.dev/cmd/vet), [sloglint](https://github.com/go-simpler/sloglint))

# Optional items:
1. [ ] Explore swagger alternatives
   - [ ] Explore [huma](https://github.com/danielgtaylor/huma) or [humagin](https://pkg.go.dev/github.com/danielgtaylor/huma/v2/adapters/humagin)
2. [ ] Explore different web frameworks
   - [ ] Explore chi
3. [ ] Add authentication

# Install migrate tool

Based on [link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

# Database migration (golang-migrate/migrate)

Install `migrate` tool using the [link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

```shell
migrate create -ext sql -dir internal/repository/db/migrations -seq <migration name in snake case>
```
