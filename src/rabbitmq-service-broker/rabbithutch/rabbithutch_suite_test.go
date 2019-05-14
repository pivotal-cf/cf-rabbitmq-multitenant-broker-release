package rabbithutch_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRabbithutch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rabbithutch Suite")
}

type fakeBody struct {
	Closed bool
}

func (f *fakeBody) Read(p []byte) (n int, err error) {
	return
}

func (f *fakeBody) Close() error {
	f.Closed = true
	return nil
}
