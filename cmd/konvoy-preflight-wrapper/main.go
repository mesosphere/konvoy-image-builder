package main

import (
	"fmt"
	"os"

	konvoywrapper "github.com/mesosphere/konvoy-image-builder/cmd/konvoy-image-wrapper/cmd"
)

func main() {
	args := []string{"check", "preflight"}
	args = append(args, os.Args[1:]...)
	err := konvoywrapper.NewRunner().Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
