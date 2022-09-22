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
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "replace-instances",
			Usage: "replace cluster instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "cluster-id",
					Value:  "",
					Usage:  "ID of target ECS cluster",
					EnvVar: "AWS_ECS_CLUSTER_ID",
				},
				cli.IntFlag{
					Name:  "max-attempt",
					Value: 40,
					Usage: "maximum attempts of waiter config",
				},
				cli.IntFlag{
					Name:  "waiter-delay",
					Value: 20,
					Usage: "delay of wiater config",
				},
			},
			Action: func(c *cli.Context) error {
				return handler.ReplaceClusterInstnces(c)
			},
		},
	}

	if os.Getenv("AWS_REGION") == "" {
		log.Println("missing region")
		os.Exit(1)
	}
	if os.Getenv("AWS_ACCESS_KEY") == "" {
		log.Println("missing access key")
		os.Exit(1)
	}
	if os.Getenv("AWS_SECRET_KEY") == "" {
		log.Println("missing secret key")
		os.Exit(1)
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
