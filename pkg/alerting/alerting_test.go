package alerting_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/access-log-monitor/pkg/alerting"
	"github.com/sjeandeaux/access-log-monitor/pkg/statistic"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alerting Suite")
}

var _ = Describe("Default Alertor", func() {
	Describe("Alerting ", func() {
		Context("With empty traffic", func() {
			It("should send nothing", func() {
				traffics := make(<-chan statistic.Traffic, 5)
				alertor := alerting.NewDefaultAlertor(5, 10, 120*time.Second)
				alerts := alertor.Alert(context.TODO(), traffics)
				Consistently(alerts).ShouldNot(Receive())
			})
		})

		Context("With high traffic", func() {
			It("should send a alert with high traffic", func() {
				traffics := make(chan statistic.Traffic, 5)
				traffics <- statistic.Traffic{Nb: 20}
				traffics <- statistic.Traffic{Nb: 25}
				traffics <- statistic.Traffic{Nb: 30}
				traffics <- statistic.Traffic{Nb: 25}
				traffics <- statistic.Traffic{Nb: 13}
				alertor := alerting.NewDefaultAlertor(5, 10, 120*time.Second)
				alerts := alertor.Alert(context.TODO(), traffics)
				//nothing because nbHits = 20 + 25 + 30 + 25 + 13 = 113 / 120 < 10
				Consistently(alerts).ShouldNot(Receive())

				traffics <- statistic.Traffic{Nb: 1107}
				alert := &alerting.Alert{}
				//triggers a alert because nbHits = 25 + 30 + 25 + 13 + 1107 = 1200 / 120 == 10
				Eventually(alerts).Should(Receive(alert))
				Ω(alert.Status).Should(Equal(alerting.HighTraffic))
				Ω(alert.Hits).Should(Equal(int64(1200)))

			})
		})

		Context("With high traffic and slowed down traffic", func() {
			It("should send a alert with high traffic and recovered alert", func() {
				traffics := make(chan statistic.Traffic, 5)
				traffics <- statistic.Traffic{Nb: 1120}
				traffics <- statistic.Traffic{Nb: 21}
				traffics <- statistic.Traffic{Nb: 21}
				traffics <- statistic.Traffic{Nb: 21}
				traffics <- statistic.Traffic{Nb: 21}
				alertor := alerting.NewDefaultAlertor(5, 10, 120*time.Second)
				alerts := alertor.Alert(context.TODO(), traffics)

				alert := &alerting.Alert{}
				Eventually(alerts).Should(Receive(alert))
				Ω(alert.Status).Should(Equal(alerting.HighTraffic))
				Ω(alert.Hits).Should(Equal(int64(1204)))

				traffics <- statistic.Traffic{Nb: 6}
				Eventually(alerts).Should(Receive(alert))
				Ω(alert.Status).Should(Equal(alerting.Recovered))
				//nbHits = 21 * 4 + 6 = 90
				Ω(alert.Hits).Should(Equal(int64(90)))
			})
		})
	})
})
