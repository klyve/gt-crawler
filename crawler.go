package main

import (
	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/GlidingTracks/gt-crawler/chrome"
	"github.com/GlidingTracks/gt-crawler/sites"
	"github.com/MarkusAJacobsen/jConfig-go"
	"github.com/Sirupsen/logrus"
	"sync"
)


func main() {
	conf, err := getConfig()
	if err != nil {
		logrus.Fatal("Could not get config, storing result locally")
		// CHECK STORAGE AND UPLOAD RESIDUALS
	}

	var wg sync.WaitGroup

	links := crawl(&wg)
	upload(conf, links, &wg)

	wg.Wait()
}

func crawl(wg *sync.WaitGroup) (links []string){
	defer wg.Done()

	c:= &chrome.Chrome{}
	cSites := []sites.ChromeSite{&sites.XContestChrome{}}
	crawlRes := make(chan []string)

	wg.Add(1)
	go c.Crawl(cSites, crawlRes)
	links = <- crawlRes
	close(crawlRes)

	return
}

func upload(conf *State, links []string, wg *sync.WaitGroup) {
	defer wg.Done()

	uploadFinished := make(chan bool)
	up := &Upload{
		Auth: auth.FAuth{},
	}

	wg.Add(1)
	go up.UploadLinks(links, uploadFinished, conf)
	uploaded := <-uploadFinished
	close(uploadFinished)


	if uploaded {
		logrus.Info("Links uploaded")
	}
}

func getConfig() (state *State, err error) {
	conf := jConfigGo.Config{}

	if err = conf.CreateConfig("state"); err != nil {
		return
	}

	state = &State{}
	err = conf.Get(state)

	return
}
