GO_VERSION := 1.17

# Common values used throughout the Makefile, not intended to be configured.
TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

.PHONY: clean
clean:
	rm -f api $(PACKAGED_TEMPLATE)

.PHONY: build
build: clean lambda

.PHONY: run
run: build
	sam local start-api --docker-network emcy-sam-backend --profile default -p 8080

api: ./lambda/api/main.go
	go build -o api ./lambda/api/main.go

.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) api

lint:
	golangci-lint run

test:
	go test ./...
