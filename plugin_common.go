package main

import (
	"crypto/md5"
	"encoding/hex"
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
		"qianxinti",
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
	case "qianxinti":
		return &PluginQianxinTi{}
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

// MD5 -
func MD5(text string) string {
	ctx := md5.New()
	_, _ = ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
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

func secThief() {
	// TODO: https://sec.thief.one/atom.xml
}

func freebuf() {
	// TODO: https://search.freebuf.com/search/?search=%E6%BC%8F%E6%B4%9E%20%E9%A2%84%E8%AD%A6#article
}

func openwall() {
	// TODO: https://www.openwall.com/lists/oss-security/
}
