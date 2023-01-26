fmt:
	@go fmt ./...

tidy:
	@go mod tidy

vendor: tidy
	@go mod vendor

gosec:
	@gosec ./...

build: fmt vendor
	@go build ./cmd/update-ami/main.go

mockgen:
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface/interface.go -destination=./internal/mock/asgiface.go -package=mocks
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go -destination=./internal/mock/ec2iface.go -package=mocks
	mockgen -source=./vendor/github.com/aws/aws-sdk-go/service/ecs/ecsiface/interface.go -destination=./internal/mock/ecsiface.go -package=mocks

test: 
	go test -cover ./... -coverprofile=cover.out
