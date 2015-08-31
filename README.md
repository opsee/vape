vape
=====

A tokin' generator and login management. It will manage teams/organizations also. We'll see.

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

# generate the swagger files
make swagger
```

