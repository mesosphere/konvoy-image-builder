# Using custom source AMIs for creating AWS Machine Images

## Prerequisites
- AWS account

## Using custom source AMIs

When using KIB for building machine images to Amazon, the default source AMIs that we provide are based on looking up an AMI based on the owner, and a filter for that operating system and version.

You can view an example of that with the provided `centos-79.yaml` snippet below:

```yaml
download_images: true

packer:
  ami_filter_name: "CentOS 7.9.2009 x86_64"
  ami_filter_owners: "125523088429"
  distribution: "CentOS"
  distribution_version: "7.9"
  source_ami: ""
  ssh_username: "centos"
  root_device_name: "/dev/sda1"
...
```

At times, a particular upstream AMI may not be available in your region, or something could be renamed, or perhaps you want to provide a custom AMI for whatever reason you need.

If this is the case, you will want to edit, or create your own, yaml file that looks up based on the `source_ami` field.

For example, [CentOS also provides an image](https://wiki.centos.org/Cloud/AWS) on the [AWS marketplace](https://aws.amazon.com/marketplace/pp/prodview-foff247vr2zfw) which you can subscribe to for free.

Once you select the source AMI that you want, you can declare that when running your build command:

```bash
konvoy-image build aws --source-ami ami-0123456789 path/to/ami/centos-79.yaml
```

Alternatively, if you want to add it to your yaml file, or are making your own file, you can do that as well.
You just need to add that AMI ID into the `source_ami` in the yaml file:

```yaml
download_images: true

packer:
  ami_filter_name: ""
  ami_filter_owners: ""
  distribution: "CentOS"
  distribution_version: "7.9"
  source_ami: "ami-123456789"
  ssh_username: "centos"
  root_device_name: "/dev/sda1"
...
```

When you're done selecting your `source_ami`, you can build your KIB image as you would normally:

```bash
konvoy-image build path/to/ami/centos-79.yaml
```
