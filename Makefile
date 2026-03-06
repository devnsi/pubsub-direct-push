default: help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

#=============================================
##@ Setup
setup: ## Setup required resources.
	@brew install protobuf || apt install protobuf-compiler
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@protoc --version
	@protoc-gen-go --version
	@protoc-gen-go-grpc --version

#=============================================
##@ Development
help-gcs: ## Show help.
	@docker run --rm fsouza/fake-gcs-server -help

stub: ## Start GCS with stub pubsub.
	@type fake-gcs-server >/dev/null 2>&1 || go install github.com/fsouza/fake-gcs-server@latest
	@PUBSUB_EMULATOR_HOST=localhost:8085 \
	fake-gcs-server \
	  -scheme http \
	  -port 4443 \
	  -backend memory \
	  -event.pubsub-project-id=test-project \
	  -event.pubsub-topic=test-topic

#=============================================
##@ CI/CD
generate: ## Generate the server.
	@protoc --go_out=internal --go-grpc_out=internal internal/receive/pubsub.proto

start: ## Build the application.
	@go run main.go

.PHONY: build
build: generate ## Build the application.
	@go build -o build/bin main.go

container: build ## Wrap the application in a container.
	@docker build -t pubsub-direct-push:latest .

install: ## Install from GitHub.
	go install github.com/devnsi/pubsub-direct-push@latest

#=============================================
##@ Verify
bucket: ## Create bucket.
	@curl -X POST http://localhost:4443/storage/v1/b -H "Content-Type: application/json" -d '{"name":"common"}'

upload: bucket ## Upload file to bucket.
	@curl -X POST -H "Content-Type: application/json" --data '{"key1":"value1","key2":"value2"}' "http://localhost:4443/upload/storage/v1/b/common/o?uploadType=media&name=o1/g1/f1.txt"
