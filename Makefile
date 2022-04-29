TARGET_FILE = dm

clean:
	rm -rf $(TARGET_FILE)

clean-test:
	@go fmt ./...
	@go clean -testcache

build-deps:
	@go mod tidy

test: clean-test
	go test -p 1 ./...

build:
	@go build -o $(TARGET_FILE)
