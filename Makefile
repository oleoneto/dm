GOBIN := $(GOPATH)/bin
IMAGE_NAME = dm
LIB_NAME = db-migrator-lib
TARGET_FILE = dm

clean:
	rm -rf $(TARGET_FILE)

clean-test:
	@go fmt ./...
	@go clean -testcache

clean-docker:
	docker rmi $(IMAGE_NAME)

build-deps:
	@go mod tidy

test: clean-test
	go test -p 1 ./...

build: build-deps test
	@go build -o $(TARGET_FILE)

build-for-docker:
	@go build -o $(TARGET_FILE)

install: build
	@go env -w GOBIN=$(GOBIN)
	@go install
	@mv $(GOBIN)/$(LIB_NAME)  $(GOBIN)/$(TARGET_FILE)

docker:
	docker build . -t $(IMAGE_NAME)

run-docker:
	@docker run --name $(IMAGE_NAME) -it --rm -p 3809:3809 -e DATABASE_URL=${DATABASE_URL} -v ${PWD}/examples:/app/migrations $(IMAGE_NAME)
