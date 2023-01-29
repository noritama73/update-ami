package services

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/noritama73/update-ami/internal/mocks"
)

func Test_WaitUntilContainerInstanceDrained(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEcsIface := mocks.NewMockECSAPI(ctrl)
	gomock.InOrder(
		mockEcsIface.EXPECT().DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(testEcsCluster),
			ContainerInstances: []*string{aws.String(testInstanceArn1)},
		}).Return(&ecs.DescribeContainerInstancesOutput{
			ContainerInstances: []*ecs.ContainerInstance{
				{
					Ec2InstanceId:        aws.String(testInstanceID1),
					Status:               aws.String(ecs.ContainerInstanceStatusActive),
					ContainerInstanceArn: aws.String(testInstanceArn1),
					RunningTasksCount:    aws.Int64(1),
				},
			},
		}, nil).Times(3),
		mockEcsIface.EXPECT().DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(testEcsCluster),
			ContainerInstances: []*string{aws.String(testInstanceArn1)},
		}).Return(&ecs.DescribeContainerInstancesOutput{
			ContainerInstances: []*ecs.ContainerInstance{
				{
					Ec2InstanceId:        aws.String(testInstanceID1),
					Status:               aws.String(ecs.ContainerInstanceStatusActive),
					ContainerInstanceArn: aws.String(testInstanceArn1),
					RunningTasksCount:    aws.Int64(0),
				},
			},
		}, nil).Times(1),
		mockEcsIface.EXPECT().DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(testEcsClusterErr),
			ContainerInstances: []*string{aws.String(testInstanceArn2)},
		}).Return(&ecs.DescribeContainerInstancesOutput{
			ContainerInstances: []*ecs.ContainerInstance{
				{
					Ec2InstanceId:        aws.String(testInstanceID2),
					Status:               aws.String(ecs.ContainerInstanceStatusActive),
					ContainerInstanceArn: aws.String(testInstanceArn2),
					RunningTasksCount:    aws.Int64(1),
				},
			},
		}, nil).Times(3),
	)

	mockEcsService := newMockEcsService(mockEcsIface)

	t.Run("Wait for Drained", func(t *testing.T) {
		instance := ClusterInstance{
			Cluster:              testEcsCluster,
			ContainerInstanceArn: testInstanceArn1,
		}
		wConfig := CustomAWSWaiterConfig{
			MaxAttempts: 4,
			Delay:       1,
		}
		assert.NoError(t, mockEcsService.WaitUntilContainerInstanceDrained(instance, wConfig))
	})

	t.Run("Over Attempts Limit", func(t *testing.T) {
		instance := ClusterInstance{
			Cluster:              testEcsClusterErr,
			ContainerInstanceArn: testInstanceArn2,
		}
		wConfig := CustomAWSWaiterConfig{
			MaxAttempts: 3,
			Delay:       1,
		}
		assert.Equal(t, fmt.Errorf("the maximum number of attempts has been reached: %d", wConfig.MaxAttempts), mockEcsService.WaitUntilContainerInstanceDrained(instance, wConfig))
	})
}

func Test_WaitUntilNewInstanceRegistered(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEcsIface := mocks.NewMockECSAPI(ctrl)
	gomock.InOrder(
		mockEcsIface.EXPECT().ListContainerInstances(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(testEcsCluster),
		}).Return(&ecs.ListContainerInstancesOutput{
			ContainerInstanceArns: []*string{aws.String(testInstanceArn1)},
		}, nil).Times(2),
		mockEcsIface.EXPECT().ListContainerInstances(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(testEcsCluster),
		}).Return(&ecs.ListContainerInstancesOutput{
			ContainerInstanceArns: []*string{aws.String(testInstanceArn1), aws.String(testInstanceArn2)},
		}, nil).Times(1),
		mockEcsIface.EXPECT().ListContainerInstances(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(testEcsClusterErr),
		}).Return(&ecs.ListContainerInstancesOutput{
			ContainerInstanceArns: []*string{aws.String(testInstanceArn1), aws.String(testInstanceArn2)},
		}, nil).Times(3),
	)

	mockEcsService := newMockEcsService(mockEcsIface)

	t.Run("Wait for Register", func(t *testing.T) {
		wConfig := CustomAWSWaiterConfig{
			MaxAttempts: 3,
			Delay:       1,
		}
		assert.NoError(t, mockEcsService.WaitUntilNewInstanceRegistered(testEcsCluster, 2, wConfig))
	})

	t.Run("Over Attempts Limit", func(t *testing.T) {
		wConfig := CustomAWSWaiterConfig{
			MaxAttempts: 3,
			Delay:       1,
		}
		assert.Equal(t,
			fmt.Errorf("the maximum number of attempts has been reached: %d", wConfig.MaxAttempts),
			mockEcsService.WaitUntilNewInstanceRegistered(testEcsClusterErr, 3, wConfig))
	})
}
