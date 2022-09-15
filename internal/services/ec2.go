package services

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/urfave/cli"
)

type EC2Service interface {
	TerinateInstance(instance ClusterInstance) error
}

type ec2Service struct {
	svc ec2iface.EC2API
}

func NewEC2Service(c *cli.Context) (EC2Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, err
	}
	return &ec2Service{
		svc: ec2.New(sess),
	}, nil
}

func (s *ec2Service) TerinateInstance(instance ClusterInstance) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instance.InstanceID)},
	}
	_, err := s.svc.TerminateInstances(input)
	return err
}

var _ EC2Service = (*ec2Service)(nil)
