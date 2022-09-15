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
	log.Println(clusterInstances)

	for _, instance := range clusterInstances {
		// インスタンスをドレイン( update-container-instances-state )
		err := ecsService.DrainContainerInstances(instance)
		if err != nil {
			log.Println(err)
		}

		// インスタンスをクラスタから外す
		err = ecsService.DeregisterContainerInstance(instance)
		if err != nil {
			log.Println(err)
		}

		// インスタンスをterminate(termiane-instance)
		err = ec2Service.TerinateInstance(instance)
		if err != nil {
			log.Println(err)
		}

		// 新しいインスタンスが登録されるのを待つ(ヘルスチェックの猶予は300秒)
		time.Sleep(300 * time.Second)

		// ecsサービスを--force-new-deployment
		err = ecsService.UpdateECSServiceByForce(instance)
		if err != nil {
			log.Println(err)
		}

		// 全てのインスタンスが更新されるまで繰り返す
		//（ひとまず最初に取得したインスタンスを全てterminateしたら正常終了？）
	}
	return nil
}
