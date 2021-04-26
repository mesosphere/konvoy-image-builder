package main

import (
	"fmt"
	"os"

	"github.com/mesosphere/konvoy-image-builder/cmd/konvoy-image-wrapper/cmd"
)

func main() {
	err := cmd.NewRunner().Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encountered: %s\n", err)
	}
}
