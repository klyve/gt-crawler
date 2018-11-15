package chrome

import (
	"context"
	"github.com/GlidingTracks/gt-crawler/sites"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/chromedp"
	"time"
)

type Chrome struct{}

func (ch Chrome) Crawl(v []sites.ChromeSite, pipe chan []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	c := CreateInstance(ctx)

	for i := range v {
		links, err := v[i].Crawl(c, ctx)
		if err != nil {
			logrus.Error("Error during crawling, err: ", err.Error())
			return
		}

		for i := range links {
			logrus.Info(links[i])
		}

		pipe <- links
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
