package log_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-parsor/pkg/log"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Entry", func() {

	Describe("Read section", func() {
		context := func(request, section string) {
			Context(fmt.Sprintf("With %q as request", request), func() {
				It(fmt.Sprintf("should return %q as section", section), func() {
					le := &log.Entry{
						Request: request}
					Ω(le.Section()).Should(Equal(section))
				})
			})
		}

		context("", "/")
		context("/", "/")
		context("/pages", "/pages")
		context("/pages/", "/pages")
		context("/pages/creation", "/pages")
		context("/pages/creation/third", "/pages")

	})

})
