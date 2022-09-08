fmt:
	@go fmt ./...

build:
	@go build

update-ami: build
	@./update-ami hello

replace: build
	@./update-ami replace-instances
