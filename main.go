package main

import (
	"os"
	"os/signal"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	err        error
	debug      bool = false
	signalChan chan os.Signal
)

func doJob() {
	for _, pn := range getPlugins() {
		p := pluginFactry(pn)
		err := p.Crawl()
		if err != nil {
			logger.Errorln(err)
			continue
		}
		err = addWarings(p.Result())
		if err != nil {
			logger.Errorln(err)
		}
	}
}

func init() {
	if os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true" {
		debug = true
		conf.Server.Debug = debug
		logger = initLogger("itm.log", logrus.DebugLevel)
		logger.Debugln("Debug Mode Running...")
	} else {
		logger = initLogger("itm.log", logrus.InfoLevel)
	}
}

func main() {
	err = loadConfig()
	if err != nil {
		logger.Errorln(err)
		return
	}
	initDatabase()

	if debug {
		doJob()
		return
	}

	c := cron.New()
	c.AddFunc("* */10 * * * *", func() {
		logger.Println("Start Job...")
		doJob()
	})
	c.Start()

	// 等待中断信号以优雅地关闭服务器
	signalChan = make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	logger.Debugln("Get Signal:", sig)
	logger.Println("Shutdown")
}
