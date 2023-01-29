package handler

import (
	"log"
	"time"

	"github.com/urfave/cli"

	"github.com/noritama73/update-ami/internal/services"
)

////////////////////////////////////////////////////////////////////////////////////
// 1.既存のコンテナインスタンスのIDを控える
// 2.ASGのdesired countを1増やす(maxを越えてしまうならそっちも一旦+1？)
// 3.新しいインスタンスが追加されるのを待つ
// 4.古いインスタンスを1つドレインする
// 5.ドレインされたらderegister→terminate
// 6.インスタンスが増えるのを待つ
// 7.サービスを強制更新
// 8.ちょっと待つ
// 9.4.に戻る
// 10.古いインスタンスをドレインし切るタイミングでdesired countを1減らす（元に戻す）
////////////////////////////////////////////////////////////////////////////////////

func ReplaceClusterInstnces(c *cli.Context) error {
	clusterName := c.String("cluster")
	asgName := c.String("asg-name")
	if c.String("asg-name") == "" {
		asgName = clusterName
	}

	ec2Service, ecsService, asgService := services.NewServices(c.String("profile"))
	log.Println("successfully initialize sessions")

	waiterConfig := services.CustomAWSWaiterConfig{
		MaxAttempts: c.Int("max-attempt"),
		Delay:       c.Int("waiter-delay"),
	}

	clusterInstances, err := ecsService.ListContainerInstances(clusterName)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range clusterInstances {
		log.Printf("Instance is found: %v", v.InstanceID)
	}
	cap, err := asgService.DescribeAutoScalingGroups(asgName)
	if err != nil {
		log.Println("couldn't describe autoscaling group")
		return err
	}
	desiredCount := cap
	log.Printf("Desired capacity: %d", cap)

	if err := asgService.UpdateDesiredCapacity(asgName, int64(desiredCount+1)); err != nil {
		log.Println("couldn't update desired capacity")
		return err
	}
	log.Println("waiting for a new instance to be registered")
	if err := ecsService.WaitUntilNewInstanceRegistered(clusterName, len(clusterInstances)+1, waiterConfig); err != nil {
		log.Println(err)
		return err
	}
	if err := ecsService.UpdateECSServiceByForce(clusterName); err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
	newcap, err := asgService.DescribeAutoScalingGroups(asgName)
	if err != nil {
		log.Println("couldn't describe autoscaling group")
		return err
	}
	log.Printf("increased desired capacity: %d", newcap)

	for i, instance := range clusterInstances {
		log.Println("**************************************************")
		log.Printf("working on: %v (%d / %d)", instance.InstanceID, i+1, len(clusterInstances))

		log.Printf("Draining: %v", instance.InstanceID)
		if err := ecsService.DrainContainerInstances(instance); err != nil {
			log.Println(err)
			continue
		}

		if err := ecsService.WaitUntilContainerInstanceDrained(instance, waiterConfig); err != nil {
			log.Println(err)
		} else {
			log.Printf("Drained: %v", instance.InstanceID)
		}

		if err := ecsService.DeregisterContainerInstance(instance); err != nil {
			log.Println(err)
		} else {
			log.Printf("Deregistered: %v", instance.InstanceID)
		}

		if err := ec2Service.TerinateInstance(instance); err != nil {
			log.Println(err)
		} else {
			log.Printf("Terminated: %v", instance.InstanceID)
		}

		if (i + 1) == len(clusterInstances) {
			break
		}

		log.Println("waiting for a new instance to be registered")
		if err := ecsService.WaitUntilNewInstanceRegistered(clusterName, len(clusterInstances)+1, waiterConfig); err != nil {
			log.Println(err)
		}

		if err := ecsService.UpdateECSServiceByForce(clusterName); err != nil {
			log.Println(err)
		}
		time.Sleep(10 * time.Second)
	}

	if err := asgService.UpdateDesiredCapacity(asgName, int64(desiredCount)); err != nil {
		log.Println("couldn't update desired capacity")
		return err
	}
	log.Println("Reseted desired capacity")

	log.Println("Success!")
	return nil
}
