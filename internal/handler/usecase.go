package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/noritama73/update-ami/internal/services"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func outputMachineImage(mi services.MachineImage, cluster string) {
	fmt.Printf("******************Image used in %s******************\n", cluster)
	fmt.Printf("Name: %s\n", mi.Name)
	fmt.Printf("Description: %s\n", mi.Description)
	fmt.Printf("ImageID: %s\n", mi.ImageID)
	fmt.Printf("Architecture: %s\n", mi.Architecture)
	fmt.Printf("PlantFormDatails: %s\n", mi.PlatformDetails)
}

func checkDifferenceBetweenAsgAndInstances(asg *autoscaling.Group) {
	for _, override := range asg.MixedInstancesPolicy.LaunchTemplate.Overrides {
		for _, i := range asg.Instances {
			if *i.InstanceType != *override.InstanceType {
				log.Printf("\x1b[33mInstanceType is different; Instance: %s, asg: %s\x1b[0m", *i.InstanceType, *override.InstanceType)
				if !checkWantToContinue() {
					os.Exit(0)
				}
			}
		}
	}
}

// 処理を続けたいか中断したいかを標準入力から受け付ける
// 中断時に後処理が必要な場合は、呼び出し元でそれを行ってからexitする
func checkWantToContinue() bool {
	fmt.Print("Continue? (y/n): ")
	var input string
	_, err := fmt.Scan(&input) // 3つの入力値を受け付ける
	if err != nil {
		panic(err)
	}
	if strings.ToLower(input) != "y" {
		log.Println("Aborted")
		return false
	}
	return true
}
