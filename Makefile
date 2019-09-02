test:
	richgo test -cover -v ./pkg/...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

update-result:
	go run ./tools/update-result ./pkg/parser/testdata/input/query ./pkg/parser/testdata/result/query
