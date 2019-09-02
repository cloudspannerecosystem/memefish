.PHONY: test
test:
	richgo test -cover -v ./pkg/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: update-result
update-result:
	richgo test -v ./pkg/parser/parser_test.go -update

.PHONY: update-mod
update-mod:
	go mod tidy
