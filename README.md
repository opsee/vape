vape
=====

A tokin' generator and login management. It will manage teams/customers also. We'll see.

### migrations

```
# install the migrate tool
go get github.com/mattes/migrate

# create a migration file
migrate -url postgres://snorecone@localhost/vape_dev?sslmode=disable -path ./migrations create my_migration_name

# migrate up
make migrate

```

### swagger

Documentation for annotating api and model files can be found here: https://github.com/yvasiyarov/swagger/wiki/Declarative-Comments-Format


```
# install the go swagger tool
go get github.com/yvasiyarov/swagger

# you'll also need a spec converter tool since we don't have swagger 2.0 in go yet
npm install -g api-spec-converter

# generate the swagger file
make swagger
```
