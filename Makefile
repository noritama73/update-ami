fmt:
	@go fmt ./...

vendor:
	@go mod vendor

gosec:
	@gosec ./...

build: fmt vendor
	@go build -o update-ami ./cmd/main.go

replace: build
	@./update-ami replace-instances

set_cred:
	@./scripts/set_aws_credentials.sh $(TOKEN_CODE)
