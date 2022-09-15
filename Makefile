fmt:
	@go fmt ./...

vendor:
	@go mod vendor

gosec:
	@gosec ./...

build: fmt vendor
	@go build -o update-ami ./cmd/main.go
