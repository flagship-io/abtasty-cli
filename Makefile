GOBIN ?= $$(go env GOPATH)/bin

doc: 
	go run docs/gen.go

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'github.com/flagship-io/abtasty-cli/cmd/version.Version=${ABTASTY_CLI_VERSION}'" -o abtasty-cli

test: SHELL:=/bin/bash

test:
	mkdir -p coverage
	go test -v -race `go list ./... | grep -v cmd/feature-experimentation/analyze | grep -v cmd/feature-experimentation/resource | grep -v docs` -coverprofile coverage/cover.out.tmp
	cat coverage/cover.out.tmp | grep -v "mock_\|cmd/feature-experimentation/analyze" | grep -v "mock_\|cmd/feature-experimentation/resource" | grep -v "mock_\|docs"> coverage/cover.out
	go tool cover -html=coverage/cover.out -o coverage/cover.html
	go tool cover -func=coverage/cover.out

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test -race `go list ./... | grep -v cmd/feature-experimentation/analyze | grep -v cmd/feature-experimentation/resource | grep -v docs` -coverprofile cover.out.tmp -covermode=atomic
	cat coverage/cover.out.tmp | grep -v "mock_\|cmd/feature-experimentation/analyze" | grep -v "mock_\|cmd/feature-experimentation/resource | grep -v "mock_\|docs" > cover.out
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml