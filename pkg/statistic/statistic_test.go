package statistic_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-monitor/pkg/log"
	"github.com/sjeandeaux/access-log-monitor/pkg/statistic"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Statistic Suite")
}

var _ = Describe("Traffic", func() {

	Describe("Adding a log entry", func() {
		Context("With a log entry nil", func() {
			It("should do nothing", func() {
				zero := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))
				stat := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))
				Ω(stat).Should(Equal(zero))
				stat.Add(nil)
				Ω(stat).Should(Equal(zero))
			})
		})

		Context("With a log entry GET and in success", func() {
			It("should add value", func() {
				stat := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))

				expected := &statistic.Traffic{
					Start:                  mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					End:                    mustParseTime("Wed, 09 May 2018 16:30:39 +0000"),
					HitBySections:          map[string]int64{"/report": 1},
					HitByUsers:             map[string]int64{"james": 1},
					HitUnauthorizedByUsers: map[string]int64{},
					HitForbiddenByUsers:    map[string]int64{},
					HitByMethods:           map[string]int64{"GET": 1},
					HitByHosts:             map[string]int64{"127.0.0.1": 1},
					Nb:                     1,
					Nb2XX:                  1,
					Nb3XX:                  0,
					Nb4XX:                  0,
					Nb5XX:                  0,
					Size:                   123,
				}
				stat.Add(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					Method:        "GET",
					Request:       "/report/test",
					Status:        int16(200),
					Size:          123,
				})
				Ω(stat).Should(Equal(expected))
			})
		})

		Context("With a log entry POST and in failure", func() {
			It("should add value", func() {
				stat := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))

				expected := &statistic.Traffic{
					Start:                  mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					End:                    mustParseTime("Wed, 09 May 2018 16:30:39 +0000"),
					HitBySections:          map[string]int64{"/report": 1},
					HitByUsers:             map[string]int64{"james": 1},
					HitUnauthorizedByUsers: map[string]int64{},
					HitForbiddenByUsers:    map[string]int64{"james": 1},
					HitByMethods:           map[string]int64{"POST": 1},
					HitByHosts:             map[string]int64{"127.0.0.1": 1},
					Nb:                     1,
					Nb2XX:                  0,
					Nb3XX:                  0,
					Nb4XX:                  1,
					Nb5XX:                  0,
					Size:                   123,
				}
				stat.Add(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					Method:        "POST",
					Request:       "/report/test",
					Status:        int16(403),
					Size:          123,
				})
				Ω(stat).Should(Equal(expected))
			})
		})

		Context("With a log entry POST and in failure but not in the slot", func() {
			It("shouldn't add value", func() {
				stat := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))
				expected := statistic.NewTraffic(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), mustParseTime("Wed, 09 May 2018 16:30:39 +0000"))
				stat.Add(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:38 +0000"),
					Method:        "POST",
					Request:       "/report/test",
					Status:        int16(403),
					Size:          123,
				})
				Ω(stat).Should(Equal(expected))
			})
		})

	})

})

var _ = Describe("Hit", func() {
	Describe("Ordering a map", func() {
		Context("With a nil", func() {
			It("should return nil", func() {
				Ω(statistic.Order(nil)).Should(BeNil())
			})
		})

		Context("With a empty map", func() {
			It("should return empty list", func() {
				Ω(statistic.Order(map[string]int64{})).Should(Equal([]statistic.Hit{}))
			})
		})

		Context("With map with values", func() {
			It("should return ordered list", func() {
				Ω(statistic.Order(map[string]int64{"two": 100, "one": 200})).Should(Equal([]statistic.Hit{{Name: "one", Hits: 200}, {Name: "two", Hits: 100}}))
			})
		})
	})
})

// Must parse the time.Time time.RFC1123 TODO avoid copy/paste
func mustParseTime(date string) time.Time {
	p, err := time.Parse(time.RFC1123, date)
	if err != nil {
		panic(err)
	}
	return p
}
