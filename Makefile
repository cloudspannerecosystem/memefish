.PHONY: test
test: bin/richgo
	@echo
	@echo "  (x x) < memefish: test"
	@echo "  /|||\\"
	@echo
	@echo go test ./pkg/...
	@bin/richgo test -cover \
	                 -coverprofile=cover.out \
	                 -covermode=count \
	                 -coverpkg=github.com/MakeNowJust/memefish/pkg/... \
	                 ./pkg/...
	@echo go build ./example/... ./tools/...
	@go build -o /dev/null ./example/... ./tools/...

.PHONY: lint
lint: bin/golangci-lint
	@echo
	@echo "  (x x) < memefish: lint"
	@echo "  /|||\\"
	@echo
	bin/golangci-lint run ./...

.PHONY: ci
ci: lint test

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: update-result
update-result: bin/richgo
	bin/richgo test -v ./pkg/parser/parser_test.go -update

.PHONY: update-mod
update-mod:
	go mod tidy

.PHONY: install-dep
install-dep: bin/richgo bin/golangci-lint

bin/richgo:
	go build -o bin/richgo github.com/kyoh86/richgo

bin/golangci-lint:
	go build -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
