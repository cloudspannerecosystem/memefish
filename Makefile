.PHONY: test
test:
	@echo
	@echo "  (x x) < memefish: test"
	@echo "  /|||\\"
	@echo
	@echo go test ./...
	@go test -cover \
	         -coverprofile=cover.out \
	         -covermode=count \
	         -coverpkg=github.com/cloudspannerecosystem/memefish/... \
	         ./...
	@echo go build ./examples/... ./tools/...
	@go build -o /dev/null ./examples/... ./tools/...

.PHONY: lint
lint: bin/golangci-lint
	@echo
	@echo "  (x x) < memefish: lint"
	@echo "  /|||\\"
	@echo
	bin/golangci-lint run ./...

.PHONY: docs
docs:
	@echo
	@echo "  (x x) < memefish: docs"
	@echo "  /|||\\"
	@echo
	cd docs && hugo

.PHONY: ci
ci: lint test

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: update-result
update-result:
	go test -v ./parser_test.go -update

.PHONY: update-mod
update-mod:
	go mod tidy

.PHONY: install-dep
install-dep: bin/golangci-lint

bin/golangci-lint:
	GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
