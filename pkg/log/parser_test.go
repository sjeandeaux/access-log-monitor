package log_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-monitor/pkg/log"
	. "github.com/sjeandeaux/access-log-monitor/pkg/matcher"
)

// Must parse the time.Time time.RFC1123
func mustParseTime(date string) time.Time {
	p, err := time.Parse(time.RFC1123, date)
	if err != nil {
		panic(err)
	}
	return p
}

var _ = Describe("Parser", func() {
	Describe("Parsing the line", func() {
		Context("With a correct line", func() {
			It("should parse it", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 123`)
				Ω(err).Should(BeNil())
				Ω(le).Should(CmpEqual(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					Method:        "GET",
					Request:       "/report",
					Status:        int16(200),
					Size:          123,
				}))

			})
		})
		Context("With a correct line without size", func() {
			It("should parse it", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 -`)
				Ω(err).Should(BeNil())
				Ω(le).Should(CmpEqual(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					Method:        "GET",
					Request:       "/report",
					Status:        int16(200),
					Size:          0,
				}))

			})
		})

		Context("With a unparseable line not enough elements", func() {
			It("should fail", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 -`)
				Ω(err).ShouldNot(BeNil())
				Ω(le).Should(BeNil())
			})
		})

		Context("With a unparseable line on date", func() {
			It("should fail", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`127.0.0.1 id james [09/May/201:16:25:39 +0000] "GET /report HTTP/1.0" 200 -`)
				Ω(err).ShouldNot(BeNil())
				Ω(le).Should(BeNil())
			})
		})

		Context("With a unparseable line on status", func() {
			It("should fail", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" nope -`)
				Ω(err).ShouldNot(BeNil())
				Ω(le).Should(BeNil())
			})
		})

		Context("With a unparseable line on size", func() {
			It("should fail", func() {
				parser := log.NewDefaultParser()
				le, err := parser.Parse(`127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 nope`)
				Ω(err).ShouldNot(BeNil())
				Ω(le).Should(BeNil())
			})
		})
	})
})
