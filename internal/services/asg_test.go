package services

import (
	"fmt"
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
	testAsgOldMaxSize  = int64(4)
	testAsgNewCapacity = int64(3)
	testAsgNewMaxSize  = int64(5)

	testAsgNameErr     = "asg-test-err"
	testAsgErrCapacity = int64(0)
)

var (
	testAsgs    []*autoscaling.Group
	testAsgsErr []*autoscaling.Group
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
		MaxSize:              aws.Int64(testAsgOldMaxSize),
	}
	testAsgs = append(testAsgs, testAsg)

	mockAsgIface := mocks.NewMockAutoScalingAPI(ctrl)
	mockAsgIface.EXPECT().DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(testAsgName)},
	}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: testAsgs,
	}, nil).Times(1)
	mockAsgIface.EXPECT().DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(testAsgNameErr)},
	}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: testAsgsErr,
	}, nil).Times(1)

	mockAsgService := newMockAsgService(mockAsgIface)

	t.Run("Describe", func(t *testing.T) {
		asg, err := mockAsgService.DescribeAutoScalingGroups(testAsgName)
		require.NoError(t, err)
		assert.Equal(t, testAsgOldCapacity, *asg.DesiredCapacity)
		assert.Equal(t, testAsgOldMaxSize, *asg.MaxSize)
	})

	t.Run("no group", func(t *testing.T) {
		_, err := mockAsgService.DescribeAutoScalingGroups(testAsgNameErr)
		assert.Equal(t, fmt.Errorf("there is no autoscaling group: %s", testAsgNameErr), err)
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
		MaxSize:              aws.Int64(testAsgNewMaxSize),
	}).Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil).Times(1)

	mockAsgService := newMockAsgService(mockAsgIface)

	t.Run("Update Capacity", func(t *testing.T) {
		assert.NoError(t, mockAsgService.UpdateAutoScalingGroup(testAsgName, testAsgNewCapacity, testAsgNewMaxSize))
	})
}
