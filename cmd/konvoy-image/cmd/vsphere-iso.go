//

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	vSphereISOExample = "vsphere-iso --vsphere-datacenter=dc1 --vsphere-cluster=zone1 --vsphere-network=Public --vsphere-datastore=datastore1 images/vsphere-iso/ubuntu-2004.yaml"
	vSphereISOUse     = "vsphere-iso <image.yaml>"
)

func checkVSphereISOFlags(cmd *cobra.Command, args []string) error {
	// credentials
	if vCenterServer, _ := cmd.Flags().GetString("vcenter-server"); vCenterServer == "" {
		if vCenterServer = os.Getenv("VSPHERE_SERVER"); vCenterServer == "" {
			return fmt.Errorf("vcenter-server or VSPHERE_SERVER required")
		}
	}

	if vSphereUser, _ := cmd.Flags().GetString("vsphere-user"); vSphereUser == "" {
		if vSphereUser = os.Getenv("VSPHERE_USERNAME"); vSphereUser == "" {
			return fmt.Errorf("vsphere-user or VSPHERE_USERNAME required")
		}
	}

	if vSpherePassword, _ := cmd.Flags().GetString("vsphere-password"); vSpherePassword == "" {
		if vSpherePassword = os.Getenv("VSPHERE_PASSWORD"); vSpherePassword == "" {
			return fmt.Errorf("vsphere-password or VSPHERE_PASSWORD required")
		}
	}

	return nil
}

func NewvSphereISOBuildCmd() *cobra.Command {
	flags := &buildCLIFlags{}
	cmd := &cobra.Command{
		Use:     vSphereISOUse,
		Short:   "build and provision vsphere images from ISO",
		Example: vSphereISOExample,
		Args:    cobra.ExactArgs(1),
		PreRunE: checkVSphereISOFlags,
		Run: func(cmd *cobra.Command, args []string) {
			runBuild(args[0], flags)
		},
	}

	initBuildvSphereISOFlags(cmd.Flags(), flags)
	return cmd
}

func NewvSphereISOGenerateCmd() *cobra.Command {
	flags := &generateCLIFlags{}
	cmd := &cobra.Command{
		Use:     vSphereISOUse,
		Short:   "generate files relating to building vsphere images from ISO",
		Example: vSphereISOExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runGenerate(args[0], flags)
		},
	}

	initGeneratevSphereISOFlags(cmd.Flags(), flags)
	return cmd
}

func initBuildvSphereISOFlags(fs *flag.FlagSet, buildFlags *buildCLIFlags) {
	initGenerateArgs(fs, &buildFlags.generateCLIFlags)
	initvSphereISOArgs(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, buildFlags)
}

func initGeneratevSphereISOFlags(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	initGenerateArgs(fs, generateFlags)
	initvSphereISOArgs(fs, generateFlags)
}

func initvSphereISOArgs(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.VSphereISO = &app.VSphereISOArgs{}
	addvSphereISOArgs(fs, gFlags.userArgs.VSphereISO)
}

func addvSphereISOArgs(fs *flag.FlagSet, vSphereISOArgs *app.VSphereISOArgs) {
	fs.StringVar(
		&vSphereISOArgs.ISOURL,
		"iso-url",
		"",
		"replace the templates iso url",
	)
	fs.StringVar(
		&vSphereISOArgs.ISOURL,
		"iso-checksum",
		"",
		"replace the templates iso checksum",
	)

	fs.StringVar(
		&vSphereISOArgs.VCenterServer,
		"vcenter-server",
		"",
		"vCenter server address (or environment variable VSPHERE_SERVER)",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereUser,
		"vsphere-user",
		"",
		"vSphere user (or environment variable VSPHERE_USER)",
	)
	fs.StringVar(
		&vSphereISOArgs.VSpherePassword,
		"vsphere-password",
		"",
		"vSphere password (or environment variable VSPHERE_PASSWORD)",
	)
	fs.BoolVar(
		&vSphereISOArgs.VSphereInsecureConnection,
		"vsphere-insecure-connection",
		false,
		"ignore SSL certificate errors",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereClusterName,
		"vsphere-cluster",
		"",
		"vSphere cluster name",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereDataCenter,
		"vsphere-datacenter",
		"",
		"vSphere data center name",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereDataStoreName,
		"vsphere-datastore",
		"",
		"vSphere datastore to be used",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereFolder,
		"vsphere-folder",
		"",
		"place VM in vSphere folder",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereResourcePool,
		"vsphere-resource-pool",
		"",
		"specify vSphere resource pool for the build VM",
	)
	fs.StringVar(
		&vSphereISOArgs.VSphereNetwork,
		"vsphere-network",
		"",
		"specify vSphere resource pool for the build VM",
	)

	fs.StringVar(
		&vSphereISOArgs.SSHUsername,
		"ssh-username",
		"",
		"specify the initial user which gets SUDO permissions",
	)
}
