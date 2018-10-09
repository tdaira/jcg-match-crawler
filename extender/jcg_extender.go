package extender

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/tdaira/gocrawl"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var rxOk = regexp.MustCompile(`^((http://sv\.j\-cg\.com/compe/view/tour/1244)|(http://sv\.j\-cg\.com/compe/view/match/\d+/\d+))$`)

type JCGExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
}

func (x *JCGExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	x.replaceOnClickURL(doc)
	matchInfo := x.getMatchInfo(doc)
	if matchInfo != nil {
		byte, err := json.Marshal(matchInfo)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile("./data/out", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		file.Write(append(byte, "\n"...))
	}

	// Return nil and true - let gocrawl find the links
	return nil, true
}

// Repalace onclick reference with href.
func (x *JCGExtender) replaceOnClickURL(doc *goquery.Document) {
	s := doc.Find("li.match[onclick*=location\\.href]")
	s.Each(func(i int, s *goquery.Selection) {
		onClickStr := s.AttrOr("onclick", "")
		onClickURL := strings.Replace(onClickStr, "location.href=", "", -1)
		onClickURL = strings.Replace(onClickURL, "'", "\"", -1)
		s.ReplaceWithHtml("<a href=" + onClickURL + "><\a>")
	})
}

func (x *JCGExtender) getMatchInfo(doc *goquery.Document) *MatchInfo {
	matchInfo := &MatchInfo{}
	s := doc.Find("div.team div.name")
	if s.Size() == 0 {
		return nil
	}
	s.Each(func(i int, s *goquery.Selection) {
		name := s.Find("p.name_abbr a").Text()
		s.Find("a[target=\"_svp\"]").Each(func(i int, s *goquery.Selection) {
			url := s.AttrOr("href", "")
			icon := s.Find("img").AttrOr("src", "")
			deck := Deck{Name: name, URL: url, DeckType: icon}
			matchInfo.DeckList = append(matchInfo.DeckList, deck)
		})
	})
	s = doc.Find("div.game ul.game_list li.clearfix")
	s.Each(func(i int, s *goquery.Selection) {
		name := s.Find("span:nth-child(3)").Text()
		deckType := s.Find("span:nth-child(4)").Text()
		win := Win{Name: name, DeckType: deckType}
		matchInfo.WinList = append(matchInfo.WinList, win)
	})
	return matchInfo
}

// Override Filter for our need.
func (x *JCGExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}
