// Package information on project
package information_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-monitor/pkg/information"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Information Suite")
}

var _ = Describe("Information", func() {

	Describe("Printing the information", func() {
		Context("With printable value", func() {
			It("should be non empty", func() {
				Ω(information.Print()).ShouldNot(BeEmpty())
			})
		})
	})
})
