package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type ASGService interface {
	DescribeAutoScalingGroup(name string) (*autoscaling.Group, error)
	UpdateAutoScalingGroup(name string, desiredCapacity, maxSize int64) error
}

type asgService struct {
	svc autoscalingiface.AutoScalingAPI
}

type DescribeAutoScalingGroupsOutput struct {
	DesiredCapacity int64
}

func (s *asgService) DescribeAutoScalingGroup(name string) (*autoscaling.Group, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(name)},
	}
	resp, err := s.svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}
	if len(resp.AutoScalingGroups) < 1 {
		return nil, fmt.Errorf("there is no autoscaling group: %s", name)
	}
	if len(resp.AutoScalingGroups) > 1 {
		return nil, fmt.Errorf("there is more than 1 autoscaling group: %s", name)
	}

	return resp.AutoScalingGroups[0], err
}

func (s *asgService) UpdateAutoScalingGroup(name string, desiredCapacity, maxSize int64) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(desiredCapacity),
		MaxSize:              aws.Int64(maxSize),
	}
	_, err := s.svc.UpdateAutoScalingGroup(input)

	return err
}

var _ ASGService = (*asgService)(nil)
