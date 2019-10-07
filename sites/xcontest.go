package sites

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/GlidingTracks/gt-crawler/crawlTime"
	jConfigGo "github.com/MarkusAJacobsen/jConfig-go"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const baseURL = "https://www.xcontest.org/world/en/flights/daily-score-pg/"
const filterDateQuery = "#filter[date]="

// XContestChrome struct used for instanced calls?
type XContestChrome struct{}

type xContestConfig struct {
	CrawledDates []string
}

// Crawl crawls the daily score records from xcontest.org
func (xcc XContestChrome) Crawl(ctx context.Context) (sl []string, err error) {
	var nodes []*cdp.Node

	url, date, _ := getURL(true)
	logrus.Error(url)

	task := getLinksFromURL(url, &nodes)

	parentCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	logrus.Info("Before run")
	err = chromedp.Run(parentCtx, task)
	logrus.Info("After run")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("Shutting down initial instance")

	visitQueue := filterRealFlightLinks(nodes)

	sl, err = visitDetailsPagesAndExtract(parentCtx, visitQueue)
	if err == nil || len(sl) == 0 {
		return
	}

	// Ensure crawled date has been written to config before launching a new run
	dateUpdated := writeCrawledDateToConfig(date)

	if !dateUpdated {
		logrus.Info("Could not update last crawled date, aborting")
		return
	}

	// If still finding links, recursively continue to next date
	csl, err := xcc.Crawl(parentCtx)
	if err != nil {
		logrus.Error("Error during crawling, err: ", err)
		return nil, err
	}
	logrus.Infof("XContest date: %v crawled", date)
	sl = append(sl, csl...)

	return
}

func writeCrawledDateToConfig(date string) (updated bool) {
	conf := &jConfigGo.Config{}
	err := conf.CreateConfig("xcontest")
	if err != nil {
		logrus.Errorf("Could not create confix xcontest %s", err)
		return false
	}

	xcc := xContestConfig{}
	if err := conf.Get(&xcc); err != nil {
		logrus.Errorf("Could not write date: %v to config", date)
		return false
	}

	xcc.CrawledDates = append(xcc.CrawledDates, date)
	if err := conf.Write(xcc); err != nil {
		logrus.Errorf("Could not write date: %v to config", date)
		return false
	}

	return true
}

func visitDetailsPagesAndExtract(ctx context.Context, urls []string) (sl []string, err error) {
	parentCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	for i := range urls {
		var nodes []*cdp.Node

		task := getSourceLink(urls[i], &nodes)
		err := chromedp.Run(parentCtx, task)
		if err != nil {
			logrus.Error(err)
		}

		for i := range nodes {
			if strings.Contains(nodes[i].AttributeValue("href"), ".igc") {
				sl = append(sl, nodes[i].AttributeValue("href"))
			}
		}
	}

	return
}

func getLinksFromURL(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#flights`),
		chromedp.WaitReady(`//a`, chromedp.BySearch),
		chromedp.Nodes(`//a`, nodes, chromedp.BySearch),
	}
}

func getSourceLink(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady(`#flight`, chromedp.BySearch),
		chromedp.Nodes(`th[class="igc"] > a`, nodes, chromedp.BySearch),
	}
}

func filterRealFlightLinks(found []*cdp.Node) (visitQueue []string) {
	pseudoHits := make([]string, len(found))

	for i := range found {
		pseudoHits[i] = found[i].AttributeValue("href")
	}

	for k := range pseudoHits {
		if strings.Contains(pseudoHits[k], "flights/detail") {
			visitQueue = append(visitQueue, pseudoHits[k])
		}
	}

	return
}

func getURL(pagination bool) (url string, date string, err error) {
	conf := &jConfigGo.Config{}
	err = conf.CreateConfig("xcontest")
	if err != nil {
		return
	}

	xcc := xContestConfig{}
	if err = conf.Get(&xcc); err != nil {
		return
	}

	if len(xcc.CrawledDates) == 0 {
		date = crawlTime.GetDateString(0)
		url = fmt.Sprintf("%s%s%s", baseURL, filterDateQuery, date)
		return
	}

	sort.Strings(xcc.CrawledDates)

	// Earliest date crawled
	date = crawlTime.FindPreviousDate(xcc.CrawledDates[0])
	url = fmt.Sprintf("%s%s%s", baseURL, filterDateQuery, date)

	return
}
