This it template that can be used to setup enterprise level application in Go

To generate docs use:
```shell
swag init
```

# TODO:
1. [X] Create endpoint to handle simple CRUD requests
   - [X] Choose the right framework to implement (Used `gin`, although `chia` also can be used)
   - [X] Define project structure
   - [X] Implement functionality
   - [X] Add swagger (accessible by path localhost:8080/swagger/index.html )
2. [ ] Create integration with Postgres database
   - [X] Choose the right library (`pgx` with `squirrel` can be used)
   - [ ] Implement functionality
   - [ ] Add database migration
3. [ ] Test it automatically
   - [ ] Add unit tests
   - [ ] Add integration tests
4. [ ] Infrastructure setup
   - [ ] Put application in Docker
   - [ ] Create `docker-compose.yaml` file for deployment
   - [ ] Create CI/CD pipeline and push to GitHub repository
