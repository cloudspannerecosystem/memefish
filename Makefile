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
	go run ./tools/update-result ./pkg/parser/testdata/input/query ./pkg/parser/testdata/result/query

.PHONY: update-mod
update-mod:
	go mod tidy
