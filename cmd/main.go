package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/noritama73/update-ami/internal/handler"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	app := cli.NewApp()
	app.Name = "Update AMI"
	app.Usage = "Replace ECS Cluster Instances for AMI Update"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "replace-instances",
			Usage: "replace instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "cluster-id",
					Value:  "",
					Usage:  "ID of target ECS cluster",
					EnvVar: "AWS_ECS_CLUSTER_ID",
				},
				cli.StringFlag{
					Name:   "profile",
					Value:  "",
					Usage:  "profile of AWS user",
					EnvVar: "AWS_PROFILE",
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
