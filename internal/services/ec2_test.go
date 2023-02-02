package services

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/noritama73/update-ami/internal/mocks"
)

const (
	testEc2InstanceID    = "i-XXXXXXXXXX"
	testEc2InstanceIDErr = "i-EEEEEEEEEE"

	testImageID    = "ami-XXXXXXXXXX"
	testImageIDErr = "ami-EEEEEEEEEE"
)

func newMockEc2Service(iface ec2iface.EC2API) EC2Service {
	return &ec2Service{svc: iface}
}

func Test_TerminateInstance(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEc2Iface := mocks.NewMockEC2API(ctrl)
	gomock.InOrder(
		mockEc2Iface.EXPECT().TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{aws.String(testEc2InstanceID)},
		}).Return(nil, nil).Times(1),
		mockEc2Iface.EXPECT().TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{aws.String(testEc2InstanceIDErr)},
		}).Return(nil, fmt.Errorf("error")).Times(1),
	)

	mockEc2Service := newMockEc2Service(mockEc2Iface)

	t.Run("Terminate Instance", func(t *testing.T) {
		assert.NoError(t, mockEc2Service.TerinateInstance(ClusterInstance{
			InstanceID: testEc2InstanceID,
		}))
	})

	t.Run("error", func(t *testing.T) {
		assert.Error(t, mockEc2Service.TerinateInstance(ClusterInstance{
			InstanceID: testEc2InstanceIDErr,
		}))
	})
}

func Test_GetImageID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEc2Iface := mocks.NewMockEC2API(ctrl)
	gomock.InOrder(
		mockEc2Iface.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(testEc2InstanceID)},
		}).Return(&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							ImageId: aws.String(testImageID),
						},
					},
				},
			},
		}, nil).Times(1),
		mockEc2Iface.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(testEc2InstanceIDErr)},
		}).Return(&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{},
		}, nil).Times(1),
		mockEc2Iface.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(testEc2InstanceIDErr)},
		}).Return(&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{},
				},
			},
		}, nil).Times(1),
	)

	mockEc2Service := newMockEc2Service(mockEc2Iface)

	t.Run("normal", func(t *testing.T) {
		imageID, err := mockEc2Service.GetImageID(ClusterInstance{InstanceID: testEc2InstanceID})
		require.NoError(t, err)
		assert.Equal(t, testImageID, imageID)
	})

	t.Run("no reservation", func(t *testing.T) {
		_, err := mockEc2Service.GetImageID(ClusterInstance{InstanceID: testEc2InstanceIDErr})
		assert.Equal(t, fmt.Errorf("there is no reservation: %s", testEc2InstanceIDErr), err)
	})

	t.Run("no instance", func(t *testing.T) {
		_, err := mockEc2Service.GetImageID(ClusterInstance{InstanceID: testEc2InstanceIDErr})
		assert.Equal(t, fmt.Errorf("there is no instance: %s", testEc2InstanceIDErr), err)
	})
}

func Test_DescribeImages(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockEc2Iface := mocks.NewMockEC2API(ctrl)
	gomock.InOrder(
		mockEc2Iface.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
			ImageIds: []*string{aws.String(testImageID)},
		}).Return(&ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				{
					Architecture:    aws.String(""),
					ImageId:         aws.String(testImageID),
					PlatformDetails: aws.String(""),
					Description:     aws.String(""),
					Name:            aws.String(""),
				},
			},
		}, nil).Times(1),
		mockEc2Iface.EXPECT().DescribeImages(&ec2.DescribeImagesInput{
			ImageIds: []*string{aws.String(testImageIDErr)},
		}).Return(&ec2.DescribeImagesOutput{
			Images: []*ec2.Image{},
		}, nil).Times(1),
	)

	mockEc2Service := newMockEc2Service(mockEc2Iface)

	t.Run("normal", func(t *testing.T) {
		image, err := mockEc2Service.DescribeImages(testImageID)
		require.NoError(t, err)
		assert.Equal(t, testImageID, image.ImageID)
	})

	t.Run("no image", func(t *testing.T) {
		_, err := mockEc2Service.DescribeImages(testImageIDErr)
		assert.Equal(t, fmt.Errorf("there is no image: %s", testImageIDErr), err)
	})
}
