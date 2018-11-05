package chrome

import (
	"context"
	"github.com/GlidingTracks/gt-crawler/sites"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/chromedp"
	"time"
)

type Chrome struct{}

func (ch Chrome) Crawl() (links []string, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	c := CreateInstance(ctx)

	links, err = sites.CrawlXContest(c, ctx)
	if err != nil {
		return
	}

	//Upload source links
	for i := range links {
		logrus.Info(links[i])
	}

	return
}

func CreateInstance(ctx context.Context) (ins *chromedp.CDP) {
	ins, err := chromedp.New(ctx, chromedp.WithErrorf(logrus.Printf))
	if err != nil {
		logrus.Fatal(err)
	}

	return ins
}