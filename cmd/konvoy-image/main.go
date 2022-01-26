package main

import (
	"fmt"

	"github.com/mesosphere/konvoy-image-builder/cmd/konvoy-image/cmd"
)

func main() {
	cmd.Execute()
	fmt.Println("Done!")
}
