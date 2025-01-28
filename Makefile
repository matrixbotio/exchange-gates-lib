.PHONY: unit-tests
unit-tests:
	go test -race -short -v --count 1 ./...

.PHONY: integration-tests
integration-tests:
	go test -race -run TestIntegration_ -v --count 1 ./...

.PHONY: generate
generate:
	GOFLAGS=-mod=mod go generate ./...
