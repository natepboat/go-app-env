format:
	gofmt -w .
test:
	go test -v ./...
coverage:
	go test ./... -cover
release:
	git tag $(version)
	git push origin $(version)