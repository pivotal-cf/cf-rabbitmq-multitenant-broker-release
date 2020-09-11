package management_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestManagementClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ManagementClient Suite")
}
