package log_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sjeandeaux/access-log-monitor/pkg/matcher"

	"github.com/sjeandeaux/access-log-monitor/pkg/log"
)

var _ = Describe("Normalizer", func() {

	Describe("Normalizer the lines", func() {
		Context("With a correct line", func() {
			It("should receive entries", func() {
				lines := emulateChannelLines(context.TODO())
				normalizer := log.NewDefaultNormalizer()

				entries, errs := normalizer.Normalize(context.TODO(), lines)

				data := &log.Entry{}
				Eventually(entries).Should(Receive(data))
				Ω(data).Should(CmpEqual(&log.Entry{
					RemoteHost:    "127.0.0.1",
					RemoteLogName: "id",
					AuthUser:      "james",
					Datetime:      mustParseTime("Wed, 09 May 2018 16:25:39 +0000"),
					Method:        "GET",
					Request:       "/report",
					Status:        int16(200),
					Size:          123,
				}))

				Eventually(errs).Should(Receive())

			})
		})

		Context("With a correct line and closed context", func() {
			It("should close the channels", func() {
				lines := emulateChannelLines(context.TODO())
				normalizer := log.NewDefaultNormalizer()
				ctx := context.Background()
				ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
				defer cancel()
				entries, errs := normalizer.Normalize(ctx, lines)
				Eventually(entries).Should(BeClosed())
				Eventually(errs).Should(BeClosed())
			})
		})
	})
})

// emulateChannelLines it emulates the channel lines and sends message inside each 500 * time.Millisecond
// one OK and another one KO
func emulateChannelLines(ctx context.Context) chan string {
	lines := make(chan string, 20)
	go func() {
		defer close(lines)
		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lines <- `127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 123`
				lines <- `127.0.0.1 id james [09/May/2018:16:25:39 +0000] "GET /report HTTP/1.0" 200 nope`
			}
		}

	}()
	return lines
}
