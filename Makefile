fmt:
	@go fmt ./...

tidy:
	@go mod tidy

vendor: tidy
	@go mod vendor

gosec: fmt
	@gosec ./...

build: fmt vendor
	@go build ./cmd/update-ami/main.go

mockgen:
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface/interface.go -destination=./internal/mocks/asgiface.go -package=mocks
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go -destination=./internal/mocks/ec2iface.go -package=mocks
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/ecs/ecsiface/interface.go -destination=./internal/mocks/ecsiface.go -package=mocks

test: fmt
	@go test -cover $(shell go list ./... | grep -v 'mocks') -coverprofile=coverage.out 2>&1 > test-report.out
