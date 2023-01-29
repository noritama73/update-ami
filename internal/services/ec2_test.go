package services

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/noritama73/update-ami/internal/mocks"
// )

// const (
// 	testEc2InstanceID = "i-XXXXXXXXXX"
// )

// func Test_Ec2TerminateInstance(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	t.Cleanup(ctrl.Finish)

// 	mockEc2Iface := mocks.NewMockEC2API(ctrl)
// 	gomock.InOrder(
// 		mockEc2Iface.EXPECT().TerminateInstances(&ec2.TerminateInstancesInput{
// 			InstanceIds: []*string{aws.String(testEc2InstanceID)},
// 		}).Return(nil).Times(1),
// 		mockEc2Iface.EXPECT().TerminateInstances(gomock.Any()).Return(fmt.Errorf("error")).Times(1),
// 	)

// 	mockEc2Service, _, _ := NewMockServices()
// 	t.Run("Terminate Instance", func(t *testing.T) {
// 		assert.NoError(t, mockEc2Service.TerinateInstance(ClusterInstance{
// 			InstanceID: testEc2InstanceID,
// 		}))
// 	})

// 	t.Run("error", func(t *testing.T) {
// 		assert.Error(t, mockEc2Service.TerinateInstance(ClusterInstance{
// 			InstanceID: "",
// 		}))
// 	})
// }
