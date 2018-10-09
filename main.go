package main

import (
	"github.com/tdaira/gocrawl"
	"github.com/tdaira/jcg-match-crawler/extender"
	"time"
)


func main() {
	opts := gocrawl.NewOptions(new(extender.JCGExtender))

	opts.RobotUserAgent = "JCGBot"
	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogAll

	opts.MaxVisits = 1000

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://sv.j-cg.com/compe/view/tour/1244")
}
