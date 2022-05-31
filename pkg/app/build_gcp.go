package app

import (
	"fmt"
	"os"
)

type GCPArgs struct {
	ProjectID string // the project ID to which the source VM belongs.
	Zone      string // the zone where the source VM is located.
	Network   string // the network in which to load image creation, should have .
}

func ensureGCP() error {
	_, ok := os.LookupEnv(GCPCredentialEnvVariable)
	if !ok {
		return fmt.Errorf(
			"failed to get client secret (GOOGLE_APPLICATION_CREDENTIALS): %w",
			ErrConfigRequired,
		)
	}

	return nil
}
