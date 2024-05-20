package cmd

import (
	"github.com/spf13/cobra"

	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	vsphereExample = "vsphere --datacenter dc1 --cluster zone1 --datastore nfs-store1 " +
		"--network public --template=d2iq-base-templates/d2iq-base-CentOS-7.9 " +
		"images/ami/centos-79.yaml"
	vsphereUse = "vsphere <image.yaml>"
)

func NewVSphereBuildCmd() *cobra.Command {
	flags := &buildCLIFlags{}
	cmd := &cobra.Command{
		Use:     vsphereUse,
		Short:   "build and provision vsphere images",
		Example: vsphereExample,
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			runBuild(args[0], flags)
		},
	}

	initBuildVSphereFlags(cmd.Flags(), flags)
	return cmd
}

func NewVSphereGenerateCmd() *cobra.Command {
	flags := &generateCLIFlags{}
	cmd := &cobra.Command{
		Use:     vsphereUse,
		Short:   "generate files relating to building vsphere images",
		Example: vsphereExample,
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			runGenerate(args[0], flags)
		},
	}

	initGenerateVSphereFlags(cmd.Flags(), flags)
	return cmd
}

func initBuildVSphereFlags(fs *flag.FlagSet, buildFlags *buildCLIFlags) {
	initGenerateArgs(fs, &buildFlags.generateCLIFlags)
	initVSphereArgs(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, buildFlags)
}

func initGenerateVSphereFlags(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	initGenerateArgs(fs, generateFlags)
	initVSphereArgs(fs, generateFlags)
}

func initVSphereArgs(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.VSphere = &app.VSphereArgs{}
	addVSphereArgs(fs, gFlags.userArgs.VSphere)
}

func addVSphereArgs(fs *flag.FlagSet, vsphereArgs *app.VSphereArgs) {
	fs.StringVar(
		&vsphereArgs.Template,
		"template",
		"",
		"Base template to be used. Can include folder. <templatename> or <folder>/<templatename>. "+
			"Required value: you can pass the template name through an override file or image definition file.",
	)
	fs.StringVar(
		&vsphereArgs.Datacenter,
		"datacenter",
		"",
		"The vSphere datacenter. "+
			"Required value: you can pass the datacenter name through an override file or image definition file.",
	)

	fs.StringVar(
		&vsphereArgs.Cluster,
		"cluster",
		"",
		"vSphere cluster to be used. Alternatively set host. "+
			"Required value: you can pass the cluster name through an override file or image definition file.",
	)

	fs.StringVar(
		&vsphereArgs.Host,
		"host",
		"",
		"vSphere host to be used. Alternatively set cluster. "+
			"Required value: you can pass the host name through an override file or image definition file.",
	)
	fs.StringVar(
		&vsphereArgs.Datastore,
		"datastore",
		"",
		"vSphere datastore used to build and store the image template. "+
			"Required value: you can pass the datastore name through an override file or image definition file.",
	)
	fs.StringVar(
		&vsphereArgs.Network,
		"network",
		"",
		"vSphere network used to build image template. "+
			"Ensure the host running the command has access to this network. "+
			"Required value: you can pass the network name through an override file or image definition file.",
	)
	fs.StringVar(
		&vsphereArgs.Folder,
		"folder",
		"",
		"vSphere folder to store the image template",
	)
	fs.StringVar(
		&vsphereArgs.ResourcePool,
		"resource-pool",
		"",
		"vSphere resource pool to be used to build image template",
	)

	fs.StringVar(
		&vsphereArgs.SSHPrivateKeyFile,
		"ssh-privatekey-file",
		"",
		"Path to ssh private key which will be used to log into the base image template",
	)

	fs.StringVar(
		&vsphereArgs.SSHPublicKey,
		"ssh-publickey",
		"",
		//nolint:lll // a long help line
		"Path to SSH public key which will be copied to the image template. Ensure to set ssh-privatekey-file or load the private key into ssh-agent",
	)

	fs.StringVar(
		&vsphereArgs.SSHUserName,
		"ssh-username",
		"",
		"username to be used with the vSphere image template",
	)
}
