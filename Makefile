API_NAME = dm-api
GOBIN := $(GOPATH)/bin
IMAGE_NAME = migration-api
LIB_NAME = db-migrator-lib
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

build: build-deps test
	@go build -o $(TARGET_FILE)

build-for-docker:
	@go build -o $(TARGET_FILE)

build-api:
	cd api; go build -o $(API_NAME); mv $(API_NAME) ../

install: build
	@go env -w GOBIN=$(GOBIN)
	@go install
	@mv $(GOBIN)/$(LIB_NAME)  $(GOBIN)/$(TARGET_FILE)

install-api: build-api
	@go env -w GOBIN=$(GOBIN)
	@go install

docker:
	docker build . -t $(IMAGE_NAME)
