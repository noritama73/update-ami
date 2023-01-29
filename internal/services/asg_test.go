package services

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/noritama73/update-ami/internal/mocks"
)

const (
	testAsgName        = "asg-test"
	testAsgOldCapacity = int64(2)
	testAsgNewCapacity = int64(3)
)

var (
	testAsgs []*autoscaling.Group
)

func newMockAsgService(iface autoscalingiface.AutoScalingAPI) ASGService {
	return &asgService{svc: iface}
}

func Test_DescribeAutoScalingGroups(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	testAsg := &autoscaling.Group{
		AutoScalingGroupName: aws.String(testAsgName),
		DesiredCapacity:      aws.Int64(testAsgOldCapacity),
	}
	testAsgs = append(testAsgs, testAsg)

	mockAsgIface := mocks.NewMockAutoScalingAPI(ctrl)
	mockAsgIface.EXPECT().DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(testAsgName)},
	}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: testAsgs,
	}, nil).Times(1)

	mockAsgService := newMockAsgService(mockAsgIface)

	t.Run("Describe", func(t *testing.T) {
		desiredCapacity, err := mockAsgService.DescribeAutoScalingGroups(testAsgName)
		require.NoError(t, err)
		assert.Equal(t, testAsgOldCapacity, desiredCapacity)
	})
}

func Test_UpdateDesiredCapacity(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockAsgIface := mocks.NewMockAutoScalingAPI(ctrl)
	mockAsgIface.EXPECT().UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(testAsgName),
		DesiredCapacity:      aws.Int64(testAsgNewCapacity),
	}).Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil).Times(1)

	mockAsgService := newMockAsgService(mockAsgIface)

	t.Run("Update Capacity", func(t *testing.T) {
		assert.NoError(t, mockAsgService.UpdateDesiredCapacity(testAsgName, testAsgNewCapacity))
	})
}
