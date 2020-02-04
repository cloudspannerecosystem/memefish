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
update-result: bin/richgo
	bin/richgo test -v ./pkg/parser/parser_test.go -update

.PHONY: update-mod
update-mod:
	go mod tidy

.PHONY: install-dep
install-dep: bin/golangci-lint

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.23.2
