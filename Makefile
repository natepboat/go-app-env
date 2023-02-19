format:
	gofmt -w .
test:
	go test -v ./...
coverage:
	go test ./... -cover

