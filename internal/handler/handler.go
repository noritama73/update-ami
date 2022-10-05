package handler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/noritama73/update-ami/internal/services"
)

func ReplaceClusterInstnces(c *cli.Context) error {
	isSkip := c.Bool("skip-abnormal-instance")

	ecsService, err := services.NewECSService(c)
	if err != nil {
		log.Println(err)
		return err
	}
	ec2Service, err := services.NewEC2Service(c)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("successfully initialize sessions")

	waiterConfig := services.CustomAWSWaiterConfig{
		MaxAttempts: c.Int("max-attempt"),
		Delay:       c.Int("waiter-delay"),
	}

	// クラスタのコンテナインスタンス一覧を取得
	clusterInstances, err := ecsService.ListContainerInstances(c.String("cluster"))
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range clusterInstances {
		log.Printf("Instance is found: %v", v.InstanceID)
	}

	if !validateContinuingFromStdin() {
		os.Exit(1)
	}

	for i, instance := range clusterInstances {
		log.Println("**************************************************")
		log.Printf("working on: %v (%d / %d)", instance.InstanceID, i+1, len(clusterInstances))

		// インスタンスをドレイン( update-container-instances-state )
		if err := ecsService.DrainContainerInstances(instance); err != nil {
			log.Println(err)
			if isSkip {
				continue
			} else {
				return err
			}
		}

		// ドレインされるまで待つ
		if err := ecsService.WaitUntilContainerInstanceDrained(instance, waiterConfig); err != nil {
			log.Println(err)
			if isSkip {
				continue
			} else {
				return err
			}
		}
		log.Printf("Drained: %v", instance.InstanceID)

		// インスタンスをクラスタから外す
		if err := ecsService.DeregisterContainerInstance(instance); err != nil {
			log.Println(err)
			if isSkip {
				continue
			} else {
				return err
			}
		}
		log.Printf("Deregistered: %v", instance.InstanceID)

		// インスタンスをterminate(termiane-instance)
		if err := ec2Service.TerinateInstance(instance); err != nil {
			log.Println(err)
			if isSkip {
				continue
			} else {
				return err
			}
		}
		log.Printf("Terminated: %v", instance.InstanceID)

		// 新しいインスタンスが登録されるのを待つ(ヘルスチェックの猶予は300秒)
		log.Println("waiting for a new instance to be registered")

		if err := ecsService.WaitUntilNewInstanceRegistered(c.String("cluster"), len(clusterInstances), waiterConfig); err != nil {
			log.Println(err)
			if !validateContinuingFromStdin() {
				os.Exit(1)
			}
		}

		// ecsサービスを--force-new-deployment
		if err := ecsService.UpdateECSServiceByForce(instance); err != nil {
			log.Println(err)
			if !validateContinuingFromStdin() {
				os.Exit(1)
			}
		}

		time.Sleep(10 * time.Second)

		// 全てのインスタンスが更新されるまで繰り返す
		//（ひとまず最初に取得したインスタンスを全てterminateしたら正常終了？）
	}
	log.Println("Success!")
	return nil
}

func validateContinuingFromStdin() bool {
	fmt.Print("Continue? If yes, type exactly \"yes\": ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	if s.Err() != nil {
		panic("error in scannig stdin")
	}
	return strings.ToLower(s.Text()) == "yes"
}
