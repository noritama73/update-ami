package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Service interface {
	TerinateInstance(instance ClusterInstance) error
}

type ec2Service struct {
	svc ec2iface.EC2API
}

func (s *ec2Service) TerinateInstance(instance ClusterInstance) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instance.InstanceID)},
	}
	_, err := s.svc.TerminateInstances(input)
	return err
}

var _ EC2Service = (*ec2Service)(nil)
