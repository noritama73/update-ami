package services

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type CustomAWSWaiterConfig struct {
	MaxAttempts int
	Delay       int
}

func (s *ecsService) WaitUntilContainerInstanceDrained(instance ClusterInstance, config CustomAWSWaiterConfig) error {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(instance.Cluster),
		ContainerInstances: []*string{aws.String(instance.ContainerInstanceArn)},
	}
	for i := 0; i < config.MaxAttempts; i++ {
		resp, err := s.svc.DescribeContainerInstances(input)
		if err != nil {
			return err
		}
		if len(resp.ContainerInstances) != 1 {
			return fmt.Errorf("expect ContainerInstances == 1, got %d", len(resp.ContainerInstances))
		}
		if *resp.ContainerInstances[0].RunningTasksCount == *aws.Int64(0) {
			return nil
		}
		log.Printf("Still %d tasks remained", *resp.ContainerInstances[0].RunningTasksCount)
		time.Sleep(time.Duration(config.Delay) * time.Second)
	}
	return fmt.Errorf("the maximum number of attempts has been reached: %d", config.MaxAttempts)
}

func (s *ecsService) WaitUntilNewInstanceRegistered(cluster string, desire int, config CustomAWSWaiterConfig) error {
	input := &ecs.ListContainerInstancesInput{
		Cluster: &cluster,
	}
	for i := 0; i < config.MaxAttempts; i++ {
		resp, err := s.svc.ListContainerInstances(input)
		if err != nil {
			return err
		}
		if len(resp.ContainerInstanceArns) == desire {
			return nil
		}
		log.Println("Still new instance isn't registerd")
		time.Sleep(time.Duration(config.Delay) * time.Second)
	}
	return fmt.Errorf("the maximum number of attempts has been reached: %d", config.MaxAttempts)
}
