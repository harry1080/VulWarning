package main

import "testing"

func TestPluginQianxinTiCrawl(t *testing.T) {
	p := &PluginQianxinTi{}
	p.Crawl()

	for _, x := range p.Result() {
		logger.Debugln(x)
	}

}
