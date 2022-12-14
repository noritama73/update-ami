package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type ECSService interface {
	ListContainerInstances(cluster string) ([]ClusterInstance, error)
	DrainContainerInstances(instance ClusterInstance) error
	DeregisterContainerInstance(instance ClusterInstance) error
	UpdateECSServiceByForce(cluster string) error
	WaitUntilContainerInstanceDrained(instance ClusterInstance, config CustomAWSWaiterConfig) error
	WaitUntilNewInstanceRegistered(cluster string, desire int, config CustomAWSWaiterConfig) error
}

type ClusterInstance struct {
	Cluster              string
	InstanceID           string
	Status               string
	ContainerInstanceArn string
}

type ecsService struct {
	svc ecsiface.ECSAPI
}

func (s *ecsService) ListContainerInstances(cluster string) ([]ClusterInstance, error) {
	result := make([]ClusterInstance, 0)
	lciInput := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}
	lciResp, err := s.svc.ListContainerInstances(lciInput)
	if err != nil {
		return result, err
	}
	if len(lciResp.ContainerInstanceArns) == 0 {
		return result, fmt.Errorf("there is no instance in %s", cluster)
	}
	dciInput := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(cluster),
		ContainerInstances: lciResp.ContainerInstanceArns,
	}
	dciResp, err := s.svc.DescribeContainerInstances(dciInput)
	if err != nil {
		return result, err
	}
	for _, v := range dciResp.ContainerInstances {
		ci := ClusterInstance{
			Cluster:              cluster,
			InstanceID:           aws.StringValue(v.Ec2InstanceId),
			Status:               aws.StringValue(v.Status),
			ContainerInstanceArn: aws.StringValue(v.ContainerInstanceArn),
		}
		result = append(result, ci)
	}
	return result, nil
}

func (s *ecsService) DrainContainerInstances(instance ClusterInstance) error {
	input := &ecs.UpdateContainerInstancesStateInput{
		Cluster:            aws.String(instance.Cluster),
		ContainerInstances: []*string{aws.String(instance.ContainerInstanceArn)},
		Status:             aws.String(ecs.ContainerInstanceStatusDraining),
	}
	_, err := s.svc.UpdateContainerInstancesState(input)
	return err
}

func (s *ecsService) DeregisterContainerInstance(instance ClusterInstance) error {
	input := &ecs.DeregisterContainerInstanceInput{
		Cluster:           aws.String(instance.Cluster),
		ContainerInstance: aws.String(instance.ContainerInstanceArn),
		Force:             aws.Bool(true),
	}
	_, err := s.svc.DeregisterContainerInstance(input)
	return err
}

func (s *ecsService) UpdateECSServiceByForce(cluster string) error {
	dsInput := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}
	dsResp, err := s.svc.ListServices(dsInput)
	if err != nil {
		return err
	}
	for _, serviceArn := range dsResp.ServiceArns {
		usInput := &ecs.UpdateServiceInput{
			Cluster:            aws.String(cluster),
			Service:            aws.String(*serviceArn),
			ForceNewDeployment: aws.Bool(true),
		}
		_, err := s.svc.UpdateService(usInput)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ ECSService = (*ecsService)(nil)
