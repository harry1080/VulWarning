package main

import "testing"

func TestAliyunCrawl(t *testing.T) {
	p := &PluginAliyun{}
	p.Crawl()
}
