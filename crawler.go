package main

import (
	"github.com/GlidingTracks/gt-crawler/chrome"
	"github.com/Sirupsen/logrus"
)

type Crawler interface {
	Crawl() ([]string, error)
}

func main() {
	c := &chrome.Chrome{}

	links, err := c.Crawl()
	if err != nil {
		logrus.Fatal(err)
	}

	finished := make(chan bool)
	go uploadLinks(links, finished)

	uploaded := <- finished

	if uploaded {
		logrus.Info("Links uploaded")
	}
}


