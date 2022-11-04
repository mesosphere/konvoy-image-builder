package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

/*
Configurations:
basic
offline
offline-fips
offline-nvidia
nvidia
fips
*/

/*
OS:

centos 7.9
rhel 7.9
rhel 8.4
rhel 8.6
sles 15 sp3
oracle 7
flatcar
ubuntu 18.04
ubuntu 20.04

infra:
AWS
Azure
OVA
GCP

*/
func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir(".."))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project

	golang := client.Container().From("golang:1.18")
	src := client.Host().Workdir()
	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/kib", src).WithWorkdir("/kib")
	path := "cmd/konvoy-image/main.go"
	out := "bin/konvoy-image"
	golang.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "build", path, "-o", out},
	})
	// get reference to build output directory in container
	output := golang.File(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, out)
	if err != nil {
		return err
	}
	return nil
}
