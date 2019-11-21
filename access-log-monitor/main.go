package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/access-log-monitor/access-log-monitor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
