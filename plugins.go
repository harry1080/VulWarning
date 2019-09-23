package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/tidwall/gjson"
)

var (
	tr *http.Transport
)

func init() {
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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

func newCustomCollector() *colly.Collector {
	var c *colly.Collector
	c = colly.NewCollector(
		colly.UserAgent("Vul Warnings Bot"),
		colly.MaxDepth(2),
		colly.Async(true),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*help.aliyun.com*",
	// 	Parallelism: 2,
	// 	RandomDelay: 5 * time.Second,
	// })

	return c
}

func aliyun() {
	c := newCustomCollector()
	c.AllowedDomains = []string{"help.aliyun.com"}
	c.OnHTML("li.y-clear", func(e *colly.HTMLElement) {
		text := e.ChildText("a[href]")
		if strings.Contains(text, "漏洞预警") {
			link := e.ChildAttr("a[href]", "href")
			link = e.Request.AbsoluteURL(link)
			_time := e.ChildText("span")
			_time = _time[:len(_time)-8]
			CreateWarings(text, link, "aliyun", "2006-01-0215:04:05", _time)
			log.Println("===========================")
			log.Println(text, _time)
			// c.Visit(link)
		}
	})

	c.OnHTML("div#se-knowledge", func(e *colly.HTMLElement) {
		// TODO: 内容处理
	})
	log.Println("Start help.aliyun.com ...")
	c.Visit("https://help.aliyun.com/noticelist/9213612.html")
	c.Wait()
}

func cert360() {
	var err error
	var req *http.Request
	// TODO: https://cert.360.cn/warning

	c := newCustomCollector()
	c.AllowedDomains = []string{"cert.360.cn"}
	// c.OnHTML("div.news-conent", func(e *colly.HTMLElement) {
	// 	// TODO: 内容处理
	// e.Text
	// })
	log.Println("Start cert.360.cn ...")
	client := &http.Client{Transport: tr}
	req, err = http.NewRequest("GET", "https://cert.360.cn/warning/searchbypage?length=10&start=0", nil)
	if err != nil {
		log.Println("[-] New Request Error : ", err.Error())
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[-] Do Request Error : ", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[-] Read Body Error : ", err.Error())
		return
	}
	defer resp.Body.Close()
	value := gjson.Get(string(body), "data")
	value.Array()
	for _, v := range value.Array() {
		link := fmt.Sprintf("%s?id=%s", "https://cert.360.cn/warning/detail", v.Get("id").String())
		_time := v.Get("add_time_str").String()
		title := v.Get("title").String()
		CreateWarings(title, link, "cert360", "2006-01-02 15:04", _time)
		log.Println("===========================")
		log.Println(title, _time)
		// c.Visit(link)
	}
	c.Wait()
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

// CertRespDate CertRespDate
type CertRespDate struct {
	ID         string `json:"id"`
	AddTimeStr string `json:"add_time_str"`
	Title      string `json:"title"`
}

// CertResp CertResp
type CertResp struct {
	Data []CertRespDate `json:"data[]"`
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
