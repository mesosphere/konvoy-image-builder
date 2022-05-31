package gcp

type Credentials struct {
	ID       string
	Secret   string
	TenantID string
}

func NewCredentials(clientID string, clientSecret string, tenantID string) (*Credentials, error) {
	return &Credentials{
		ID:       clientID,
		Secret:   clientSecret,
		TenantID: tenantID,
	}, nil
}
