package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/sjeandeaux/access-log-parsor/access-log-parsor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
