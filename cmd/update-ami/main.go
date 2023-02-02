package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/noritama73/update-ami/internal/handler"
)

func main() {
	app := cli.NewApp()
	app.Name = "Update AMI"
	app.Usage = "Replace ECS Cluster Instances for AMI Update"

	app.Commands = []cli.Command{
		{
			Name:  "replace-instances",
			Usage: "replace cluster instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "cluster",
					Value:    "",
					Usage:    "Name of target ECS cluster",
					EnvVar:   "AWS_ECS_CLUSTER",
					Required: true,
				},
				cli.StringFlag{
					Name:     "region",
					Value:    "",
					Usage:    "AWS region",
					EnvVar:   "AWS_REGION",
					Required: true,
				},
				cli.StringFlag{
					Name:     "profile",
					Value:    "",
					Usage:    "AWS profile",
					EnvVar:   "AWS_PROFILE",
					Required: true,
				},
				cli.IntFlag{
					Name:  "max-attempt",
					Value: 40,
					Usage: "maximum attempts of waiter config",
				},
				cli.IntFlag{
					Name:  "waiter-delay",
					Value: 20,
					Usage: "delay of waiter config",
				},
				cli.StringFlag{
					Name:  "asg-name",
					Value: "",
					EnvVar: "AWS_ASG_NAME",
					Usage: "associated asg: if not set, this will be the same value as cluster",
				},
			},
			Action: func(c *cli.Context) error {
				return handler.ReplaceClusterInstnces(c)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
