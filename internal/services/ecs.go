package services

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type ECSService interface {
}

type ClusterInstance struct {
	Cluster            string
	InstanceID         string
	Status             string
	ClusterInstanceArn string
}

type ecsService struct {
	svc ecsiface.ECSAPI
}

func NewECSService() ECSService {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("_CM_AWS_REGION"))})
	if err != nil {
		log.Fatalln("Cannot initialize session")
	}
	return &ecsService{
		svc: ecs.New(sess),
	}
}

func (s *ecsService)ListContainerInstances(cluster string) ([]ClusterInstance, error) {
	result := make([]ClusterInstance, 0)
	return result, nil
}
