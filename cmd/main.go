package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"

	"github.com/noritama73/update-ami/internal/handler"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Update AMI"
	app.Usage = "Replace ECS Cluster Instances for AMI Update"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "replace-instaces",
			Usage: "do replace",
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
					Usage:  "profile for credential",
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
