package handler

import (
	"fmt"

	"github.com/noritama73/update-ami/internal/services"
)

func outputMachineImage(mi services.MachineImage, cluster string) {
	fmt.Printf("******************Image used in %s******************\n", cluster)
	fmt.Printf("Name: %s\n", mi.Name)
	fmt.Printf("Description: %s\n", mi.Description)
	fmt.Printf("ImageID: %s\n", mi.ImageID)
	fmt.Printf("Architecture: %s\n", mi.Architecture)
	fmt.Printf("PlantFormDatails: %s\n", mi.PlatformDetails)
}
