package app

type AmazonArgs struct {
	// AMI options
	SourceAMI        string   `json:"source_ami"`
	AMIFilterName    string   `json:"ami_filter_name"`
	AMIFilterOwner   string   `json:"ami_filter_owner"`
	AWSBuilderRegion string   `json:"aws_region"`
	AMIRegions       []string `json:"ami_regions"`
	AWSInstanceType  string   `json:"aws_instance_type"`

	AMIUsers  []string `json:"ami_users"`
	AMIGroups []string `json:"ami_groups"`
}
