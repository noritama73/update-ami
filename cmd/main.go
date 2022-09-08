package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	clusterID string
)

func main() {

	app := cli.NewApp()
	app.Name = "Update AMI"
	app.Usage = "Replace ECS Cluster Instances for AMI Update"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cluster-id",
			Value:       "",
			Usage:       "ID of target ECS cluster",
			Destination: &clusterID,
			EnvVar:      "AWS_ECS_CLUSTER_ID",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "replace-instaces",
			Usage: "do replace",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "cluster-id",
					Value:       "",
					Usage:       "ID of target ECS cluster",
					Destination: &clusterID,
					EnvVar:      "AWS_ECS_CLUSTER_ID",
				},
			},
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}
