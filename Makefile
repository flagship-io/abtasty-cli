build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'github.com/flagship-io/abtasty-cli/cmd/version.Version=${ABTASTY_CLI_VERSION}'" -o abtasty-cli

test: SHELL:=/bin/bash
test:
	mkdir -p coverage
	go test -v -race `go list ./... | grep -v cmd/feature-experimentation/analyze | grep -v cmd/feature-experimentation/resource` -coverprofile coverage/cover.out.tmp
	cat coverage/cover.out.tmp | grep -v "mock_\|cmd/feature-experimentation/analyze" | grep -v "mock_\|cmd/feature-experimentation/resource" > coverage/cover.out
	go tool cover -html=coverage/cover.out -o coverage/cover.html
	go tool cover -func=coverage/cover.out