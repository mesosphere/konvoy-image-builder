package ansible

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

var _ = Describe("Provision", func() {
	BeforeEach(func() {
		// use the mock data with the container for e2e tests
		appansible.PlaybookPath = "testdata/ansible"
		app.AnsibleRunsDirectory = "generated/ansible-runs"
	})

	It("runs provision with the given inventory.yaml", func() {
		Provision("inventory.yaml")
	})

	It("runs the playbook with the given hostname", func() {
		Provision("localhost,")
	})
})

func Provision(inventory string) {
	err := app.Provision(inventory, app.ProvisionFlags{
		ExtraVars: []string{"@extra-vars.yaml"},
		RootFlags: app.RootFlags{
			Verbosity: logging.TraceLevel,
		},
	})
	Expect(err).ToNot(HaveOccurred())
}
