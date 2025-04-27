This it template that can be used to setup enterprise level application in Go

To generate docs use:
```shell
swag init
```

Starting Postgres:
```shell
docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:17-alpine
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
3. [ ] Setup logging
   - [ ] Choose library for logging or understand how to implement it
4. [ ] Test it automatically
   - [ ] Add unit tests
   - [ ] Add integration tests
   - [ ] Add test coverage report
5. [ ] Infrastructure setup
   - [ ] Put application in Docker
   - [ ] Create `docker-compose.yaml` file for deployment
   - [ ] Create CI/CD pipeline and push to GitHub repository

# Optional items:
1. [ ] Explore swagger alternatives
   - [ ] Explore [huma](https://github.com/danielgtaylor/huma) or [humagin](https://pkg.go.dev/github.com/danielgtaylor/huma/v2/adapters/humagin)
2. [ ] Explore different web frameworks
   - [ ] Explore chi

# Install migrate tool

Based on [link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

# Database migration (golang-migrate/migrate)

Install `migrate` tool using the [link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

```shell
migrate create -ext sql -dir internal/repository/db/migrations -seq <migration name in snake case>
```
