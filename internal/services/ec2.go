package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Service interface {
	TerinateInstance(instance ClusterInstance) error

	GetImageID(instance ClusterInstance) (string, error)
	DescribeImages(id string) (MachineImage, error)
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

func (s *ec2Service) GetImageID(instance ClusterInstance) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instance.InstanceID)},
	}
	res, err := s.svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}
	if len(res.Reservations) < 1 {
		return "", fmt.Errorf("there is no reservation: %s", instance.InstanceID)
	}
	if len(res.Reservations[0].Instances) < 1 {
		return "", fmt.Errorf("there is no instance: %s", instance.InstanceID)
	}

	return *res.Reservations[0].Instances[0].ImageId, nil
}

type MachineImage struct {
	Architecture    string
	ImageID         string
	PlatformDetails string
	Description     string
	Name            string
}

func (s *ec2Service) DescribeImages(id string) (MachineImage, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(id)},
	}
	res, err := s.svc.DescribeImages(input)
	if len(res.Images) < 1 {
		return MachineImage{}, fmt.Errorf("there is no image: %s", id)
	}
	image := res.Images[0]
	machineImage := MachineImage{
		Architecture:    *image.Architecture,
		ImageID:         *image.ImageId,
		PlatformDetails: *image.PlatformDetails,
		Description:     *image.Description,
		Name:            *image.Name,
	}

	return machineImage, err
}

var _ EC2Service = (*ec2Service)(nil)
