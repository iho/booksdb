# BooksDB

## Run server
```
docker-compose up
```
## Run inmemory tests
```
make s
```
## Run tests wich spins up database in Docker via ory/dockertest
```
make test
```
## Run linter
1. Install linter https://golangci-lint.run/usage/install/#local-installation
2. Run autoformater and linter
```
go get github.com/daixiang0/gci
go get mvdan.cc/gofumpt
make check
```
