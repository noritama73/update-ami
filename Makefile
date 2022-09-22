fmt:
	@go fmt ./...

vendor:
	@go mod vendor

gosec:
	@gosec ./...

build: fmt vendor
	@go build ./cmd/update-ami/main.go
