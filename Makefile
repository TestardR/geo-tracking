NAME=geo-tracking
VERSION="develop"

.PHONY: dep fmt test clean vet

deps:
	@echo "Installing dependencies ..."
	@go mod tidy && go mod download
	@echo "Installing dependencies, done!"

fmt:
	@echo "Formatting code ..."
	@go fmt ./...
	@echo "Formatting code, done!"

test:
	@echo "Running tests ..."
	@ginkgo -r -cover ./internal
	@echo "Running tests, done!"

build: deps clean
	@echo "Building ...."
	@mkdir -p ./build
	@go build -o ./build/$(NAME)
	@echo "Build done"

install:
	@echo "Installing ..."
	@go install
	@echo "Installed"

clean:
	@echo "Cleaning ..."
	@go clean && rm -rf ./build
	@echo "Cleaning done"

vet:
	@echo "Running vet ..."
	@go vet ./...
	@echo "Running vet, done!"