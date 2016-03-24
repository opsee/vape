docfile = src/github.com/opsee/vape/api/docs.go
APPENV := testenv
PROJECT := vape
REV ?= latest

all: build

clean:
	rm -fr target bin pkg

fmt:
	@gofmt -w ./

migrate:
	migrate -url $(POSTGRES_CONN) -path ./migrations up

deps:
	docker-compose up -d

swagger:
	@mkdir -p swagger
	GOPATH="$(PWD):$(PWD)/vendor:$(GOPATH)" swagger -apiPackage github.com/opsee/vape -mainApiFile github.com/opsee/vape/api/api.go -format swagger -output swagger
	@for path in authenticate signups users ; do \
		mv swagger/$$path/index.json swagger/$$path.json && rmdir swagger/$$path ; \
	done
	@echo "package api\n\nvar swaggerJson=\`" > $(docfile)
	@api-spec-converter swagger/index.json --from=swagger_1 --to=swagger_2 >> $(docfile)
	@echo "\`" >> $(docfile)
	@rm -fr swagger

build: deps $(APPENV)
	docker run \
		--link $(PROJECT)_postgres_1:postgres \
		--env-file ./$(APPENV) \
		-e "TARGETS=linux/amd64" \
		-e PROJECT=github.com/opsee/$(PROJECT) \
		-v `pwd`:/gopath/src/github.com/opsee/$(PROJECT) \
		quay.io/opsee/build-go:16
	docker build -t quay.io/opsee/$(PROJECT):$(REV) .

run: deps $(APPENV)
	docker run \
		--link $(PROJECT)_postgres_1:postgres \
		--env-file ./$(APPENV) \
		-e AWS_DEFAULT_REGION \
		-e AWS_ACCESS_KEY_ID \
		-e AWS_SECRET_ACCESS_KEY \
		-p 8081:8081 \
		-p 9091:9091 \
		--rm \
		quay.io/opsee/$(PROJECT):$(REV)

.PHONY: migrate clean all swagger docker run
