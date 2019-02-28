package config_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

func fixture(name string) string {
	path, err := filepath.Abs(filepath.Join("fixtures", name))
	Expect(err).NotTo(HaveOccurred())
	return path
}
