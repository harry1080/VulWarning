package main

import (
	"github.com/sirupsen/logrus"
)

var (
	logger = initLogger("itm.log", logrus.DebugLevel)
	err    error
)

func doJob() {
	for _, pn := range getPlugins() {
		p := pluginFactry(pn)
		if p.Crawl() == nil {
			// TODO: Save res and  pusher
			addWarings(p.Result())
		}
	}
}

func main() {
	initDatabase()

	// c := cron.New()
	// c.AddFunc("*/10 * * * *", func() {
	// 	doJob()
	// })
	// c.Start()

	doJob()

}
