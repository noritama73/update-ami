package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/urfave/cli"
)

func NewSessions(c *cli.Context) (EC2Service, ECSService) {
	opt := session.Options{
		Config:                  *aws.NewConfig(),
		Profile:                 c.String("profile"),
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		AssumeRoleDuration:      3600,
		SharedConfigState:       session.SharedConfigEnable,
	}
	sess := session.Must(session.NewSessionWithOptions(opt))
	return &ec2Service{svc: ec2.New(sess)}, &ecsService{ecs.New(sess)}
}
