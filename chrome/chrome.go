package chrome

import (
	"context"
	"github.com/GlidingTracks/gt-crawler/sites"
	"github.com/Sirupsen/logrus"
)

type Chrome struct{}

func (ch Chrome) Crawl(ctx context.Context, v []sites.ChromeSite, pipe chan []string) {
	for i := range v {

		links, err := v[i].Crawl(ctx)
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
