package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type ASGService interface {
	DescribeAutoScalingGroups(name string) (int64, error)
	UpdateDesiredCapacity(name string, newCapacity int64) error
}

type asgService struct {
	svc autoscalingiface.AutoScalingAPI
}

type DescribeAutoScalingGroupsOutput struct {
	DesiredCapacity int64
}

func (s *asgService) DescribeAutoScalingGroups(name string) (int64, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(name)},
	}
	resp, err := s.svc.DescribeAutoScalingGroups(input)

	return *resp.AutoScalingGroups[0].DesiredCapacity, err
}

func (s *asgService) UpdateDesiredCapacity(name string, newCapacity int64) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(newCapacity),
	}
	_, err := s.svc.UpdateAutoScalingGroup(input)

	return err
}

var _ ASGService = (*asgService)(nil)
