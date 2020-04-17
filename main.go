package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	logger     = initLogger("itm.log", logrus.InfoLevel)
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
	flag.BoolVar(&debug, "debug", false, "this help")
}

func main() {
	flag.Parse()
	err = loadConfig()
	if err != nil {
		logger.Errorln(err)
		return
	}

	if debug {
		logger.Debugln("Debug Mode Running...")
		conf.Server.Debug = debug
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
	logger.Println("Get Signal:", sig)
	logger.Println("Shutdown")
}
