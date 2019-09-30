package main

import (
	"context"
	"sync"
	"time"

	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/GlidingTracks/gt-crawler/chrome"
	"github.com/GlidingTracks/gt-crawler/sites"
	jConfigGo "github.com/MarkusAJacobsen/jConfig-go"
	"github.com/Sirupsen/logrus"
)

func init() {
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	logrus.Info("Starting crawler")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	conf, err := getConfig()
	if err != nil {
		logrus.Fatal("Could not get config, storing result locally")
		// CHECK STORAGE AND UPLOAD RESIDUALS
	}

	var wg sync.WaitGroup

	links := crawl(ctx, &wg)
	upload(ctx, conf, links)

	wg.Wait()
}

func crawl(ctx context.Context, wg *sync.WaitGroup) (links []string) {
	defer wg.Done()

	c := &chrome.Chrome{}
	cSites := []sites.ChromeSite{&sites.XContestChrome{}}
	crawlRes := make(chan []string)

	wg.Add(1)
	go c.Crawl(ctx, cSites, crawlRes)
	links = <-crawlRes
	close(crawlRes)

	return
}

func upload(ctx context.Context, conf State, links []string) {
	if len(links) == 0 {
		logrus.Info("No links uploaded, empty input")
		return
	}

	up := &Upload{
		Auth: auth.FAuth{},
	}

	uploaded := up.UploadLinks(ctx, links, conf)
	if uploaded {
		logrus.Info("Links uploaded")
	}
}

func getConfig() (state State, err error) {
	conf := jConfigGo.Config{}

	if err = conf.CreateConfig("state"); err != nil {
		return
	}

	state = State{}
	err = conf.Get(&state)

	return
}
