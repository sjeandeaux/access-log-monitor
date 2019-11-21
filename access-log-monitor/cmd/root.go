package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "access-log-parsor",
	Short: "Monitoring the log",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFormatter(&log.TextFormatter{})

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", log.InfoLevel.String(), "The log level")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(watchCmd)

}

func initConfig() {
	if l, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	viper.AutomaticEnv()
}
