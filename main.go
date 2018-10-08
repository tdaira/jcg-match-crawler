package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tdaira/gocrawl"
	"net/http"
	"regexp"
	"strings"
	"time"
	"fmt"
)

var rxOk = regexp.MustCompile(`^((http://sv\.j\-cg\.com/compe/view/tour/1244)|(http://sv\.j\-cg\.com/compe/view/match/\d+/\d+))$`)

type JCGExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
}

func (x *JCGExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	// Repalace onclick reference with href.
	s := doc.Find("li.match[onclick*=location\\.href]")
	s.Each(func(i int, s *goquery.Selection) {
		onClickStr := s.AttrOr("onclick", "")
		onClickURL := strings.Replace(onClickStr, "location.href=", "", -1)
		onClickURL = strings.Replace(onClickURL, "'", "\"", -1)
		fmt.Println("url: " + onClickURL)
		s.ReplaceWithHtml("<a href=" + onClickURL + "><\a>")
	})

	// Return nil and true - let gocrawl find the links
	return nil, true
}

// Override Filter for our need.
func (x *JCGExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func main() {
	opts := gocrawl.NewOptions(new(JCGExtender))

	opts.RobotUserAgent = "JCGBot"
	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogAll

	opts.MaxVisits = 1000

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://sv.j-cg.com/compe/view/tour/1244")
}
