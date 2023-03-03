package app

type VSphereArgs struct {
	// VSphere options
	Template          string `json:"template"`
	Cluster           string `json:"cluster"`
	Host              string `json:"host"`
	Datacenter        string `json:"datacenter"`
	Datastore         string `json:"datastore"`
	Network           string `json:"network"`
	Folder            string `json:"folder"`
	ResourcePool      string `json:"resource_pool"`
	SSHPrivateKeyFile string `json:"ssh_private_key_file"`
	SSHPublicKey      string `json:"ssh_public_key"`
	SSHUserName       string `json:"ssh_username"`
}
