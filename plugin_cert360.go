package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

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

// PluginCert360 -
type PluginCert360 struct {
	c   *colly.Collector
	res []*Warings
}

// Result -
func (p *PluginCert360) Result() []*Warings {
	return p.res
}

// Crawl -
func (p *PluginCert360) Crawl() error {
	params := make(map[string]string)
	params["length"] = "10"
	params["start"] = "0"
	body, err := httpGet("https://cert.360.cn/warning/searchbypage", params)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	obj := gjson.GetBytes(body, "data")
	if obj.IsArray() {
		for _, v := range obj.Array() {
			_time := v.Get("add_time_str").String()
			title := v.Get("title").String()
			desc := v.Get("description").String()
			id := v.Get("id").String()
			link := fmt.Sprintf("%s?id=%s", "https://cert.360.cn/warning/detail", id)
			p.res = append(p.res, &Warings{
				Title:    title,
				Link:     link,
				Index:    id,
				From:     "cert360",
				Desc:     desc,
				Time:     getTime("2006-01-02 15:04", _time),
				CreateAt: time.Now(),
			})
			logger.Debugln("Crwaled [Cert360]", title, _time)
		}
	}
	return nil
}
