vape
=====

A tokin' generator and login management. It will manage teams/organizations also. We'll see.


authentications
- password
- refresh
- oauth callback

logins
- delete
- update
- show

signups
- create
- list
- show
- approve
- claim



### migrations

```
# install the migrate tool
go get github.com/mattes/migrate

# create a migration file
migrate -url postgres://snorecone@localhost/vape_dev?sslmode=disable -path ./migrations create my_migration_name

# migrate up
migrate -url postgres://snorecone@localhost/vape_dev?sslmode=disable -path ./migrations up

```
