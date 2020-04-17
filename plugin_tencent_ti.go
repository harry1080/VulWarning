package main

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// PluginTencentTi -
type PluginTencentTi struct {
	c   *colly.Collector
	res []*Warings
}

// Result -
func (p *PluginTencentTi) Result() []*Warings {
	return p.res
}

// Crawl -
func (p *PluginTencentTi) Crawl() error {
	p.c = newCustomCollector([]string{"security.tencent.com"})

	p.c.OnRequest(func(r *colly.Request) {
		logger.Println("Crawling [TencentTi]", r.URL)
	})

	p.c.OnHTML("div.user_body", func(e *colly.HTMLElement) {
		title := e.ChildText("h2.body_title")
		logger.Debugln(title)
		_time := e.ChildText("div.content_rightblock > p.content_time > span")
		desc := ""
		e.ForEach("div.body_block-detail", func(i int, ex *colly.HTMLElement) {
			if strings.Contains(ex.ChildText("h3"), "更新标题") {
				desc = ex.ChildText("div")
			}
		})

		p.res = append(p.res, &Warings{
			Title:    title,
			Link:     e.Request.URL.String(),
			From:     "tencent_ti",
			Desc:     desc,
			Time:     getTime("2006-01-02 15:04:05", _time),
			CreateAt: time.Now(),
		})
	})

	p.c.OnHTML("tbody > tr > td", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.ChildAttr("a[href]", "href"))
		if strings.Contains(link, "update_detail") {
			// logger.Debugln(link)
			p.c.Visit(link)
		}
	})
	p.c.Visit("https://security.tencent.com/ti?type=1")
	p.c.Wait()

	return nil
}
