package cmd

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_mountFileEnv(t *testing.T) {
	cases := []struct {
		name       string
		envName    string
		createFile bool
		mountPath  string
		expect     func(*GomegaWithT, *Runner, error)
	}{
		{
			name:       "successfully mount file assigned to env variable",
			envName:    "TEST_VALID_PATH",
			createFile: true,
			mountPath:  "",
			expect: func(g *GomegaWithT, r *Runner, err error) {
				g.Expect(err).To(BeNil())
				g.Expect(len(r.volumes)).To(Equal(1))

				g.Expect(r.env["TEST_VALID_PATH"]).To(Equal(os.Getenv("TEST_VALID_PATH")))

				expectedMountPath := os.Getenv("TEST_VALID_PATH")
				g.Expect(r.volumes[0].source).To(Equal(expectedMountPath))
				g.Expect(r.volumes[0].target).To(Equal(expectedMountPath))
			},
		},
		{
			name:       "successfully mount file at given path",
			envName:    "TEST_ALT_MOUNT_PATH",
			createFile: true,
			mountPath:  "/test/path",
			expect: func(g *GomegaWithT, r *Runner, err error) {
				g.Expect(err).To(BeNil())
				g.Expect(len(r.volumes)).To(Equal(1))

				g.Expect(r.env["TEST_ALT_MOUNT_PATH"]).To(Equal("/test/path"))

				expectedMountPath := os.Getenv("TEST_ALT_MOUNT_PATH")
				g.Expect(r.volumes[0].source).To(Equal(expectedMountPath))
				g.Expect(r.volumes[0].target).To(Equal("/test/path"))
			},
		},
		{
			name:       "ignore if env variable is not set",
			envName:    "TEST_IGNORE_ENV_VAR",
			createFile: false,
			mountPath:  "",
			expect: func(g *GomegaWithT, r *Runner, err error) {
				g.Expect(err).To(BeNil())
				g.Expect(len(r.volumes)).To(Equal(0))
				g.Expect(os.Getenv("TEST_IGNORE_ENV_VAR")).To(Equal(""))
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			r := NewRunner()
			if c.createFile {
				f, err := os.CreateTemp("", "example")
				g.Expect(err).To(BeNil())

				t.Setenv(c.envName, f.Name())
				defer os.Remove(f.Name())
			}

			err := r.mountFileEnv(c.envName, c.mountPath)
			c.expect(g, r, err)
		})
	}

	// test invalid path
	t.Run("error on invalid path", func(t *testing.T) {
		g := NewGomegaWithT(t)
		t.Setenv("TEST_INVALID_PATH", "/path/doesnot/exists")
		r := NewRunner()
		err := r.mountFileEnv("TEST_INVALID_PATH", "")
		g.Expect(err).ToNot(BeNil())
	})

	// test path set to dir
	t.Run("error on path set to dir", func(t *testing.T) {
		g := NewGomegaWithT(t)
		t.Setenv("TEST_PATH_TO_DIR", os.TempDir())
		r := NewRunner()
		err := r.mountFileEnv("TEST_PATH_TO_DIR", "")
		g.Expect(err).ToNot(BeNil())
	})
}
