all: fmt build

build:
	gb build

clean:
	rm -fr target bin pkg

fmt:
	@gofmt -w ./

migrate:
	migrate -url $(POSTGRES_CONN) -path ./migrations up

swagger: src/github.com/opsee/vape/**/*.go
	GOPATH="$(PWD):$(PWD)/vendor:$(GOPATH)" swagger -apiPackage github.com/opsee/vape -mainApiFile github.com/opsee/vape/api/api.go -format swagger -output swagger

.PHONY: migrate clean all
