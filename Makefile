.PHONY:pre-lint
pre-lint:
	command -v gofumpt >/dev/null 2>&1 || go install mvdan.cc/gofumpt@v0.9.2
	command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

.PHONY:lint
lint: gci pre-lint
	go mod tidy
	gofumpt -l -w .
	go vet ./...
	golangci-lint run

.PHONY:pre-gci
pre-gci:
	command -v gci >/dev/null 2>&1 || go install github.com/daixiang0/gci@v0.13.7

.PHONY:gci
gci: pre-gci
	gci write --skip-generated -s standard -s default .

.PHONY:pre-api-doc
pre-api-doc:
	command -v swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@v1.16.6
