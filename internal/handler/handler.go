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

	for _, instance := range clusterInstances {
		// インスタンスをドレイン( update-container-instances-state )
		if err := ecsService.DrainContainerInstances(instance); err != nil {
			log.Println(err)
		}

		// ドレインされるまで待つ
		config := services.CustomAWSWaiterConfig{
			MaxAttempts: 40,
			Delay:       10,
		}
		if err := ecsService.WaitUntilContainerInstanceDrained(instance, config); err != nil {
			log.Println(err)
		}

		// インスタンスをクラスタから外す
		if err := ecsService.DeregisterContainerInstance(instance); err != nil {
			log.Println(err)
		}

		// インスタンスをterminate(termiane-instance)
		if err := ec2Service.TerinateInstance(instance); err != nil {
			log.Println(err)
		}

		// 新しいインスタンスが登録されるのを待つ(ヘルスチェックの猶予は300秒)
		time.Sleep(300 * time.Second)

		// ecsサービスを--force-new-deployment
		if err := ecsService.UpdateECSServiceByForce(instance); err != nil {
			log.Println(err)
		}

		// 全てのインスタンスが更新されるまで繰り返す
		//（ひとまず最初に取得したインスタンスを全てterminateしたら正常終了？）
	}
	return nil
}
