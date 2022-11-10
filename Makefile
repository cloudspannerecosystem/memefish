.PHONY: test
test:
	@echo
	@echo "  (x x) < memefish: test"
	@echo "  /|||\\"
	@echo
	@echo go test ./pkg/...
	@go test -cover \
	         -coverprofile=cover.out \
	         -covermode=count \
	         -coverpkg=github.com/MakeNowJust/memefish/pkg/... \
	         ./pkg/...
	@echo go build ./example/... ./tools/...
	@go build -o /dev/null ./example/... ./tools/...

.PHONY: lint
lint: golangci-lint
	@echo
	@echo "  (x x) < memefish: lint"
	@echo "  /|||\\"
	@echo
	golangci-lint run ./...

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
update-result: bin/richgo
	bin/richgo test -v ./pkg/parser/parser_test.go -update

.PHONY: update-mod
update-mod:
	go mod tidy

.PHONY: install-dep
install-dep: golangci-lint

golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
