package services

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/noritama73/update-ami/internal/mocks"
)

const (
	testEcsCluster   = "test-cluster"
	testInstanceArn1 = "arn:aws:ec2:region:account-id:instance/instance-id1"
	testInstanceArn2 = "arn:aws:ec2:region:account-id:instance/instance-id2"
	testInstanceID1  = "i-XXXXXXXXXX"
	testInstanceID2  = "i-YYYYYYYYYY"

	testEcsClusterErr = "cluster-error"

	testEcsserviceArn1 = "arn:aws:ecs:region:account-ad:service/my-service1"
	testEcsserviceArn2 = "arn:aws:ecs:region:account-ad:service/my-service2"
)

func newMockEcsService(iface ecsiface.ECSAPI) ECSService {
	return &ecsService{svc: iface}
}

func Test_ListContainerInstances(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEcsIface := mocks.NewMockECSAPI(ctrl)
	gomock.InOrder(
		mockEcsIface.EXPECT().ListContainerInstances(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(testEcsCluster),
		}).Return(&ecs.ListContainerInstancesOutput{
			ContainerInstanceArns: []*string{aws.String(testInstanceArn1), aws.String(testInstanceArn2)},
		}, nil).Times(1),
		mockEcsIface.EXPECT().DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(testEcsCluster),
			ContainerInstances: []*string{aws.String(testInstanceArn1), aws.String(testInstanceArn2)},
		}).Return(&ecs.DescribeContainerInstancesOutput{
			ContainerInstances: []*ecs.ContainerInstance{
				{
					Ec2InstanceId:        aws.String(testInstanceID1),
					Status:               aws.String(ecs.ContainerInstanceStatusActive),
					ContainerInstanceArn: aws.String(testInstanceArn1),
				},
				{
					Ec2InstanceId:        aws.String(testInstanceID2),
					Status:               aws.String(ecs.ContainerInstanceStatusActive),
					ContainerInstanceArn: aws.String(testInstanceArn2),
				},
			},
		}, nil).Times(1),
		mockEcsIface.EXPECT().ListContainerInstances(&ecs.ListContainerInstancesInput{
			Cluster: aws.String(testEcsClusterErr),
		}).Return(&ecs.ListContainerInstancesOutput{
			ContainerInstanceArns: []*string{},
		}, fmt.Errorf("there is no instance in %s", testEcsClusterErr)).Times(1),
	)

	mockEcsService := newMockEcsService(mockEcsIface)

	t.Run("List Container Instances", func(t *testing.T) {
		clusterInstances, err := mockEcsService.ListContainerInstances(testEcsCluster)
		require.NoError(t, err)
		assert.Equal(t, []ClusterInstance{
			{
				Cluster:              testEcsCluster,
				InstanceID:           testInstanceID1,
				Status:               ecs.ContainerInstanceStatusActive,
				ContainerInstanceArn: testInstanceArn1,
			},
			{
				Cluster:              testEcsCluster,
				InstanceID:           testInstanceID2,
				Status:               ecs.ContainerInstanceStatusActive,
				ContainerInstanceArn: testInstanceArn2,
			},
		}, clusterInstances)
	})

	t.Run("No Instance", func(t *testing.T) {
		clusterInstances, err := mockEcsService.ListContainerInstances(testEcsClusterErr)
		assert.Equal(t, fmt.Errorf("there is no instance in %s", testEcsClusterErr), err)
		assert.Len(t, clusterInstances, 0)
	})
}

func Test_DrainContainerInstances(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEcsIface := mocks.NewMockECSAPI(ctrl)
	gomock.InOrder(
		mockEcsIface.EXPECT().UpdateContainerInstancesState(&ecs.UpdateContainerInstancesStateInput{
			Cluster:            aws.String(testEcsCluster),
			ContainerInstances: []*string{aws.String(testInstanceArn1)},
			Status:             aws.String(ecs.ContainerInstanceStatusDraining),
		}).Return(nil, nil).Times(1),
	)

	mockEcsService := newMockEcsService(mockEcsIface)

	t.Run("Drain Instance", func(t *testing.T) {
		require.NoError(t, mockEcsService.DrainContainerInstances(ClusterInstance{
			Cluster:              testEcsCluster,
			ContainerInstanceArn: testInstanceArn1,
		}))
	})
}

func Test_UpdateECSServiceByForce(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEcsIface := mocks.NewMockECSAPI(ctrl)
	gomock.InOrder(
		mockEcsIface.EXPECT().ListServices(&ecs.ListServicesInput{
			Cluster: aws.String(testEcsCluster),
		}).Return(&ecs.ListServicesOutput{
			ServiceArns: []*string{aws.String(testEcsserviceArn1), aws.String(testEcsserviceArn2)},
		}, nil).Times(1),
		mockEcsIface.EXPECT().UpdateService(&ecs.UpdateServiceInput{
			Cluster:            aws.String(testEcsCluster),
			Service:            aws.String(testEcsserviceArn1),
			ForceNewDeployment: aws.Bool(true),
		}).Return(&ecs.UpdateServiceOutput{}, nil).Times(1),
		mockEcsIface.EXPECT().UpdateService(&ecs.UpdateServiceInput{
			Cluster:            aws.String(testEcsCluster),
			Service:            aws.String(testEcsserviceArn2),
			ForceNewDeployment: aws.Bool(true),
		}).Return(&ecs.UpdateServiceOutput{}, nil).Times(1),
	)

	mockEcsService := newMockEcsService(mockEcsIface)

	t.Run("Update Services", func(t *testing.T) {
		require.NoError(t, mockEcsService.UpdateECSServiceByForce(testEcsCluster))
	})
}
