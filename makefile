.PHONY: unittest
unittest:
	@ginkgo ./...

.PHONY: run
run:
	@go run ./cmd/proxy/main.go

.PHONY: build
build:
	@go build -o bin/proxy ./cmd/proxy/main.go