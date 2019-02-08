package integrationtests_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func TestIntegrationtests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integrationtests Suite")
}

var session *gexec.Session

var _ = BeforeSuite(func() {
	pathToServiceBroker, err := gexec.Build("rabbitmq-service-broker")
	Expect(err).NotTo(HaveOccurred())

	path, err := filepath.Abs(filepath.Join("fixtures", "config.yml"))
	Expect(err).ToNot(HaveOccurred())

	command := exec.Command(pathToServiceBroker, "-configPath", path)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session.Out).Should(gbytes.Say("RabbitMQ Service Broker listening on port"))
})

var _ = AfterSuite(func() {
	session.Kill().Wait()
	gexec.CleanupBuildArtifacts()
})
