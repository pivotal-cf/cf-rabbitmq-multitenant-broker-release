package integrationtests_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	baseURL  = "http://localhost:8902/v2/"
	username = "p1-rabbit"
	password = "p1-rabbit-testpwd"
)

var (
	session   *gexec.Session
	rmqClient *rabbithole.Client
)

func TestIntegrationtests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integrationtests Suite")
}

var _ = BeforeSuite(func() {
	pathToServiceBroker, err := gexec.Build("rabbitmq-service-broker")
	Expect(err).NotTo(HaveOccurred())

	path, err := filepath.Abs(filepath.Join("fixtures", "config.yml"))
	Expect(err).ToNot(HaveOccurred())

	command := exec.Command(pathToServiceBroker, "-configPath", path, "-port", "8902")
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(brokerIsServing).Should(BeTrue())

	rmqClient, err = rabbithole.NewClient("http://127.0.0.1:15672", "guest", "guest")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Check TCP connections for any CLOSE_WAIT connections
	// 3D38 is hex for 15672 (3D37 is 15671), a two digit number next to it means it's in the remote_address column
	// 08 is the index of CLOSE_WAIT in the tcp_states enum. See https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/net/tcp_states.h
	cmd := "cat /proc/net/tcp | grep -e ':3D38 08 ' -e ':3D37 08 ' | wc -l"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error checking for tcp connections in CLOSE_WAIT state: ", err)
	}
	number, _ := strconv.Atoi(strings.TrimSuffix(string(out), "\n"))
	Expect(number).To(Equal(0))

	session.Kill().Wait()
	gexec.CleanupBuildArtifacts()
})

func doRequest(method, url string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, url, body)
	Expect(err).NotTo(HaveOccurred())

	req.SetBasicAuth(username, password)
	req.Header.Set("X-Broker-API-Version", "2.14")

	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())

	bodyContent, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	Expect(resp.Body.Close()).To(Succeed())
	return resp, bodyContent
}

func brokerIsServing() bool {
	_, err := http.Get(baseURL)
	return err == nil
}
