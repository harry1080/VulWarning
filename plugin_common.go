package main

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

// Plugin -
type Plugin interface {
	Crawl() error
	Result() []*Warings
}

func getPlugins() []string {
	return []string{
		"aliyun",
		"cert360",
		"tencentti",
	}
}

func pluginFactry(name string) Plugin {
	switch name {
	case "aliyun":
		return &PluginAliyun{}
	case "cert360":
		return &PluginCert360{}
	case "tencentti":
		return &PluginTencentTi{}
	default:
		return nil
	}
}

func newCustomCollector(domains []string) *colly.Collector {
	var c *colly.Collector
	c = colly.NewCollector(
		colly.UserAgent("Vul Warnings Bot"),
		colly.MaxDepth(2),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
	)
	c.AllowedDomains = domains

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	return c
}

func getTime(_timeFormat, _time string) time.Time {
	t, err := time.Parse(_timeFormat, _time)
	if err != nil {
		log.Println(err.Error())
		t = time.Now()
	}
	return t
}

func getFeature(content string) {
	// TODO: getFeature
	/*
		timestamp
		cve
		cvss
		vul_type
		version
		product
		summary
		ref
	*/
}

func qianxinTi() {
	// TODO: https://ti.qianxin.com/advisory/
}

func secThief() {
	// TODO: https://sec.thief.one/atom.xml
}

func cnnvd() {
	// TODO: http://www.cnnvd.org.cn/web/cnnvdnotice/querylist.tag
}

func anquanke() {
	// TODO: https://www.anquanke.com/tag/%E6%BC%8F%E6%B4%9E%E9%A2%84%E8%AD%A6
}

func freebuf() {
	// TODO: https://search.freebuf.com/search/?search=%E6%BC%8F%E6%B4%9E%20%E9%A2%84%E8%AD%A6#article
}

func openwall() {
	// TODO: https://www.openwall.com/lists/oss-security/
}
