package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sjeandeaux/access-log-parsor/pkg/alerting"
	"github.com/sjeandeaux/access-log-parsor/pkg/display"
	"github.com/sjeandeaux/access-log-parsor/pkg/file"
	"github.com/sjeandeaux/access-log-parsor/pkg/log"
	"github.com/sjeandeaux/access-log-parsor/pkg/statistic"
	"github.com/spf13/cobra"
)

const layoutTime = "02/Jan/2006:15:04:05 -0700"

type commandLine struct {
	file  string
	start string

	aggregationDuration time.Duration

	alertDuration  time.Duration
	alertThreshold int64
}

func (c *commandLine) Check() error {
	if c.aggregationDuration > c.alertDuration {
		return fmt.Errorf("The alerting uses the traffic aggregation, aggregation-duration should be inferior than alert-duration")
	}

	if c.alertDuration%c.aggregationDuration != 0 {
		return fmt.Errorf("The alerting uses the traffic aggregation, aggregation-duration should be multiple of alert-duration")
	}

	return nil

}

var cmdLine = &commandLine{}

var watchCmd = &cobra.Command{
	Use:   `watch`,
	Short: "watch an access log file",
	PreRun: func(createCmd *cobra.Command, args []string) {
		if err := cmdLine.Check(); err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}
	},
	Run: func(createCmd *cobra.Command, args []string) {

		ctx := context.Background()
		//Create the tailer on the file
		tailer, err := file.NewDefaultTailer(cmdLine.file)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		//Listen the channels
		//Create the UI
		terminalUI, err := display.NewTerminalUI()
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		//Read the start time
		start, err := time.Parse(layoutTime, cmdLine.start)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}
		//Create the normalizer
		normalizer := log.NewDefaultNormalizer()
		//Create the aggregator
		aggregator := statistic.NewDefaultAggregator(start, cmdLine.aggregationDuration)
		alertor := alerting.NewDefaultAlertor(int64(cmdLine.alertDuration.Seconds()/cmdLine.aggregationDuration.Seconds()), float64(cmdLine.alertThreshold), cmdLine.alertDuration)

		//Tail the file
		lines, _ := tailer.Tail(ctx)
		//Normalize the lines
		entries, errChan := normalizer.Normalize(ctx, lines)
		//Aggregate the entries
		traffics := aggregator.Aggregate(ctx, entries)

		forAlerting := make(chan statistic.Traffic, 20)
		forUI := make(chan statistic.Traffic, 20)

		//multi plexing, it uses the traffic channel for the UI and for the alerting
		go func() {
			for t := range traffics {
				forAlerting <- t
				forUI <- t
			}
		}()

		//Listen the entries to generate alerts
		alerts := alertor.Alert(ctx, forAlerting)

		if err := terminalUI.Run(ctx, forUI, alerts, errChan); err != nil {
			logrus.Fatal(err)
			os.Exit(1)

		}

	},
}

func init() {
	watchCmd.Flags().StringVarP(&cmdLine.file, "file", "f", "/tmp/access.log", "The access log file")
	watchCmd.Flags().StringVarP(&cmdLine.start, "start", "s", time.Now().Format(layoutTime), "The start time '02/Jan/2006:15:04:05 -0700' default is now.")
	watchCmd.Flags().DurationVarP(&cmdLine.aggregationDuration, "aggregation-duration", "", 10*time.Second, "The aggregation duration")

	watchCmd.Flags().Int64VarP(&cmdLine.alertThreshold, "alert-threshold", "", 10, "The number of requests")
	watchCmd.Flags().DurationVarP(&cmdLine.alertDuration, "alert-duration", "", 2*time.Minute, "The alert duration")
}
