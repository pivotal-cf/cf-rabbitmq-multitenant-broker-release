package integrationtests_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const url = "http://localhost:8901/v2/catalog"

var _ = Describe("/v2/catalog", func() {
	When("no credentials are provided", func() {
		It("fails with HTTP 401", func() {
			response, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})
