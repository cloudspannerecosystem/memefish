.PHONY: test
test: gen
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
lint: gen bin/golangci-lint
	@echo
	@echo "  (x x) < memefish: lint"
	@echo "  /|||\\"
	@echo
	bin/golangci-lint run ./...

.PHONY: gen
gen:
	@echo
	@echo "  (x x) < memefish: gen"
	@echo "  /|||\\"
	@echo
	go generate ./...

.PHONY: check-gen
check-gen: gen
	git diff --exit-code

.PHONY: docs
docs:
	@echo
	@echo "  (x x) < memefish: docs"
	@echo "  /|||\\"
	@echo
	cd docs && hugo mod get -u && hugo

.PHONY: ci
ci: check-gen lint test

.PHONY: fmt
fmt: gen
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
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(CURDIR)/bin v2.1.6
