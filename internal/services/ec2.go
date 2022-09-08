package services

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Service interface {
}

type ec2Service struct {
	svc ec2iface.EC2API
}

func NewEC2Service() ECSService {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("_CM_AWS_REGION"))})
	if err != nil {
		log.Fatalln("Cannot initialize session")
	}
	return &ec2Service{
		svc: ec2.New(sess),
	}
}
