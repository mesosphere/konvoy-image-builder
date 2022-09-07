# Create OS images for GCP


## Prerequisites
- Install gcloud CLI: [Instructions](https://cloud.google.com/sdk/docs/install)
- Authentication:

    Packer plugin for Google Compute requires authenticate with GCP using service account credentials. Instructions to create service account and credentials file can be found [here](https://www.packer.io/plugins/builders/googlecompute#running-outside-of-google-cloud)

    **Using CLI:**
    ```bash
    export USER=<SERVICE_ACCOUNT_USER>
    export GCP_PROJECT=<GCP_PROJECT_NAME>
    export GOOGLE_APPLICATION_CREDENTIALS=$HOME/.gcloud/credentials.json
    gcloud iam service-accounts create "$USER"
    gcloud projects add-iam-policy-binding $GCP_PROJECT --member="serviceAccount:$USER@$GCP_PROJECT.iam.gserviceaccount.com" --role=roles/compute.instanceAdmin.v1
    gcloud projects add-iam-policy-binding $GCP_PROJECT --member="serviceAccount:$USER@$GCP_PROJECT.iam.gserviceaccount.com" --role=roles/iam.serviceAccountUser
    gcloud iam service-accounts keys create $GOOGLE_APPLICATION_CREDENTIALS --iam-account="$USER@$GCP_PROJECT.iam.gserviceaccount.com"
    ```
- Network:

    A network with a firewall rule set to allow SSH traffic must be created to allow Packer to communicate to the VM provisioning image.
    **Using CLI:**
    ```shell
    export NETWORK_NAME=<NAME_OF_NETWORK>
    gcloud compute networks create "${NETWORK_NAME}" --project="$GCP_PROJECT" --subnet-mode=auto --mtu=1460 --bgp-routing-mode=regional

    gcloud compute firewall-rules create "${NETWORK_NAME}-allow-ssh" --project="$GCP_PROJECT" --network="projects/$GCP_PROJECT/global/networks/${CLUSTER_NAME}" --description=Allows\ TCP\ connections\ from\ any\ source\ to\ any\ instance\ on\ the\ network\ using\ port\ 22. --direction=INGRESS --priority=65534 --source-ranges=0.0.0.0/0 --action=ALLOW --rules=tcp:22
    ```
- Environment variables
Make sure to create a file with credentials for the service account using the instructions above.

```bash
export GOOGLE_APPLICATION_CREDENTIALS=$HOME/.gcloud/credentials.json
```

**Packer variables for GCP:**

Add the following configuration in the `image.yaml`
Substitue following variables needed for building images in GCP
- PROJECT_NAME
- ZONE
- NETWORK_NAME

```yaml
packer:
  # The source image to use to create the new image from. source_image = `distribution`-`distribution_version`
  distribution: "centos"
  distribution_version: "7-9-v20220519"
  source_image: "centos-7-v20220519" # GCP maintained public base image
  distribution_family: "centos-7"
  # The username to connect to SSH with. Required if using SSH.
  ssh_username: "centos"
  project_id: "<PROJECT_NAME>"
  zone: "<GCP_ZONE>" # List all zones using  gcloud compute zones list
  network: "<NETWORK>" # name of the network

build_name: "centos-7"
packer_builder_type: "googlecompute"
python_path: ""
```

## Create image on GCP

```bash
konvoy-image build gcp path/to/image.yaml
```

Checkout example image configurations at the [`<project_root>/images/gcp/`](../../images/gcp) directory.
