package statistic_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-monitor/pkg/log"
	. "github.com/sjeandeaux/access-log-monitor/pkg/matcher"
	"github.com/sjeandeaux/access-log-monitor/pkg/statistic"
)

var _ = Describe("Aggregator", func() {

	Describe("Aggregate the entries", func() {
		Context("With two slot", func() {
			It("should receive two slots", func() {
				entries := emulateChannelEntries(context.TODO())
				aggregator := statistic.NewDefaultAggregator(mustParseTime("Wed, 09 May 2018 16:25:39 +0000"), 3*time.Second)

				traffics := aggregator.Aggregate(context.TODO(), entries)

				data := &statistic.Traffic{}
				Eventually(traffics, 9*time.Second).Should(Receive(data))
				Ω(data).Should(CmpEqual(&statistic.Traffic{
					Start:                  mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					End:                    mustParseTime("Wed, 09 May 2018 16:25:42 +0000"),
					HitBySections:          map[string]int64{"/report": 3},
					HitByUsers:             map[string]int64{"james": 3},
					HitUnauthorizedByUsers: map[string]int64{},
					HitForbiddenByUsers:    map[string]int64{},
					HitByMethods:           map[string]int64{"GET": 3},
					HitByHosts:             map[string]int64{"127.0.0.1": 3},
					Nb:                     3,
					Nb2XX:                  3,
					Size:                   369,
				}))

			})
		})
	})
})

// emulateChannelLines it emulates the channel log.Entry
// one OK and another one KO
func emulateChannelEntries(ctx context.Context) chan log.Entry {
	entries := make(chan log.Entry, 20)
	go func() {
		defer close(entries)
		entries <- log.Entry{
			RemoteHost:    "127.0.0.1",
			RemoteLogName: "id",
			AuthUser:      "james",
			Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
			Method:        "GET",
			Request:       "/report",
			Status:        int16(200),
			Size:          123,
		}
		entries <- log.Entry{
			RemoteHost:    "127.0.0.1",
			RemoteLogName: "id",
			AuthUser:      "james",
			Datetime:      mustParseTime("Wed, 09 May 2018 16:25:41 +0000"),
			Method:        "GET",
			Request:       "/report",
			Status:        int16(200),
			Size:          123,
		}

		entries <- log.Entry{
			RemoteHost:    "127.0.0.1",
			RemoteLogName: "id",
			AuthUser:      "james",
			Datetime:      mustParseTime("Wed, 09 May 2018 16:25:42 +0000"),
			Method:        "GET",
			Request:       "/report",
			Status:        int16(200),
			Size:          123,
		}
	}()
	return entries
}
