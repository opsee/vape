docfile = src/github.com/opsee/vape/api/docs.go

all: fmt build

build:
	gb build

clean:
	rm -fr target bin pkg

fmt:
	@gofmt -w ./

migrate:
	migrate -url $(POSTGRES_CONN) -path ./migrations up

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

docker: fmt
	docker run -e POSTGRES_CONN="postgres://postgres@postgresql/vape_test?sslmode=disable" \
		--link postgresql:postgresql \
		-e "TARGETS=linux/amd64" \
		-v `pwd`:/build quay.io/opsee/build-go \
		&& docker build -t quay.io/opsee/vape .

run: docker
	docker run -e POSTGRES_CONN="postgres://postgres@postgresql/vape_test?sslmode=disable" \
		--link postgresql:postgresql \
		-e MANDRILL_API_KEY=$(MANDRILL_API_KEY) \
		-p 8081:8081 \
		-p 9091:9091 \
		--rm \
		quay.io/opsee/vape

.PHONY: migrate clean all swagger docker run
