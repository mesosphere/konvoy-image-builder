package app

type VSphereISOArgs struct {
	// AMI options
	ISOURL        string   `json:"iso_url"`
	ISOChecksum        string   `json:"iso_checksum"`
	// DistroName string `json:"iso_checksum"`
	// DistroVersion string `json:"iso_checksum"`

	VCenterServer string   `json:"vcenter_server"`
	VSphereUser string   `json:"vsphere_user"`
	VSpherePassword string   `json:"vsphere_password"`
	VSphereInsecureConnection bool   `json:"insecure_connection"`

	VSphereClusterName string   `json:"cluster"`
	VSphereDataStoreName string   `json:"datastore"`
	VSphereFolder string   `json:"folder"`
	VSphereResourcePool string   `json:"resource_pool"`
	VSphereDataCenter string   `json:"datacenter"`
	VSphereNetwork string   `json:"network"`


	SSHUsername string   `json:"ssh_username"`
}
