.PHONY: unit-tests
unit-tests:
	go test -short -v --count 1 ./...

.PHONY: integration-tests
integration-tests:
	go test -run TestIntegration_ -v --count 1 ./...
