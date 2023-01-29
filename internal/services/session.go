package services

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func NewServices(profile string) (EC2Service, ECSService, ASGService) {
	opt := session.Options{
		Config:                  *aws.NewConfig(),
		Profile:                 profile,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		AssumeRoleDuration:      3600 * time.Second,
		SharedConfigState:       session.SharedConfigEnable,
	}
	sess := session.Must(session.NewSessionWithOptions(opt))
	return &ec2Service{svc: ec2.New(sess)}, &ecsService{svc: ecs.New(sess)}, &asgService{svc: autoscaling.New(sess)}
}
