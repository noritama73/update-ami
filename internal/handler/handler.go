package handler

import (
	"log"
	"time"

	"github.com/noritama73/update-ami/internal/services"
	"github.com/urfave/cli"
)

func ReplaceClusterInstnces(c *cli.Context) error {
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

	// クラスタのコンテナインスタンス一覧を取得
	clusterInstances, err := ecsService.ListContainerInstances(c.String("cluster-id"))
	if err != nil {
		log.Println(err)
	}
	for _, v := range clusterInstances {
		log.Printf("Instance ID: %v", v.InstanceID)
	}

	for i, instance := range clusterInstances {
		log.Println("**************************************************************")
		log.Printf("working on: %v (%d / %d)", instance.InstanceID, i, len(clusterInstances))

		// インスタンスをドレイン( update-container-instances-state )
		if err := ecsService.DrainContainerInstances(instance); err != nil {
			log.Println(err)
		}

		// ドレインされるまで待つ
		config := services.CustomAWSWaiterConfig{
			MaxAttempts: 40,
			Delay:       20,
		}
		if err := ecsService.WaitUntilContainerInstanceDrained(instance, config); err != nil {
			log.Println(err)
		}
		log.Printf("Drained: %v", instance.InstanceID)

		// インスタンスをクラスタから外す
		if err := ecsService.DeregisterContainerInstance(instance); err != nil {
			log.Println(err)
		}
		log.Printf("Deregistered: %v", instance.InstanceID)

		// インスタンスをterminate(termiane-instance)
		if err := ec2Service.TerinateInstance(instance); err != nil {
			log.Println(err)
		}
		log.Printf("Terminated: %v", instance.InstanceID)

		// 新しいインスタンスが登録されるのを待つ(ヘルスチェックの猶予は300秒)
		log.Println("waiting for a new instence to be registered")
		// time.Sleep(300 * time.Second)
		if err := ecsService.WaitUntilNewInstanceRegistered(c.String("cluster-id"), len(clusterInstances), config); err != nil {
			log.Println(err)
		}

		// ecsサービスを--force-new-deployment
		if err := ecsService.UpdateECSServiceByForce(instance); err != nil {
			log.Println(err)
		}

		time.Sleep(10 * time.Second)

		// 全てのインスタンスが更新されるまで繰り返す
		//（ひとまず最初に取得したインスタンスを全てterminateしたら正常終了？）
	}
	log.Println("Success!")
	return nil
}
