package main

import (
	"encoding/json"
	"github.com/GlidingTracks/gt-crawler/chrome"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
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

	conf, err := getConfig()
	if err != nil {
		logrus.Error("Could not get config, storing result locally")
		// TODO fallback routine for saving crawled result and try to upload again at start-up if existing stored files exist
	}

	finished := make(chan bool)
	go uploadLinks(links, finished, conf.FirebaseCredentials, conf.CrawlerUID)

	uploaded := <-finished

	if uploaded {
		logrus.Info("Links uploaded")
	}
}

func getConfig() (state *State, err error) {
	confFile, err := os.OpenFile("state.json", os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer confFile.Close()

	ct, err := ioutil.ReadAll(confFile)
	if err != nil {
		return
	}

	state = &State{}
	json.Unmarshal(ct, &state)

	return
}
