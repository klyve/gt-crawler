package sites

import (
	"context"
	"github.com/GlidingTracks/gt-crawler/crawlTime"
	"github.com/MarkusAJacobsen/jConfig-go"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"sort"
	"strings"
	"time"
)

const baseUrl = "https://www.xcontest.org/world/en/flights/daily-score-pg/"
const filterDateQuery = "#filter[date]="
const pagination = "@flights[start]="

type XContestChrome struct{}

type xContestConfig struct {
	CrawledDates []string
}

// Crawl crawls the daily score records from xcontest.org
func (xcc XContestChrome) Crawl(ctx context.Context) (sl []string, err error) {
	var nodes []*cdp.Node

	url, date, _ := getURL(true)

	task := getLinksFromUrl(url, &nodes)

	ins := CreateInstance(ctx)
	err = ins.Run(ctx, task)
	if err != nil {
		return
	}

	err = ins.Shutdown(ctx)
	if err != nil {
		return
	}

	err = ins.Wait()
	if err != nil {
		return
	}

	logrus.Info("Shutting down initial instance")

	visitQueue := filterRealFlightLinks(nodes)

	sl, err = visitDetailsPagesAndExtract(visitQueue, ctx)
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
	csl, err := xcc.Crawl(ctx)
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
	conf.CreateConfig("xcontest")

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

func visitDetailsPagesAndExtract(urls []string, ctx context.Context) (sl []string, err error) {
	ins, err := chromedp.New(ctx, chromedp.WithErrorf(logrus.Printf))
	if err != nil {
		return
	}

	for i := range urls {
		var nodes []*cdp.Node

		task := getSourceLink(urls[i], &nodes)
		err := ins.Run(ctx, task)
		if err != nil {
			logrus.Error(err)
		}

		for i := range nodes {
			if strings.Contains(nodes[i].AttributeValue("href"), ".igc") {
				sl = append(sl, nodes[i].AttributeValue("href"))
			}
		}
	}

	err = ins.Shutdown(ctx)
	if err != nil {
		return
	}

	err = ins.Wait()

	return
}

func getLinksFromUrl(url string, nodes *[]*cdp.Node) chromedp.Tasks {
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
	conf.CreateConfig("xcontest")

	xcc := xContestConfig{}
	if err = conf.Get(&xcc); err != nil {
		return
	}

	if len(xcc.CrawledDates) == 0 {
		date = crawlTime.GetDateString(0)
		url = baseUrl + filterDateQuery + date
		return
	}

	sort.Strings(xcc.CrawledDates)

	// Earliest date crawled
	date = crawlTime.FindPreviousDate(xcc.CrawledDates[0])
	url = baseUrl + filterDateQuery + date

	return
}

func CreateInstance(ctx context.Context) (ins *chromedp.CDP) {
	ins, err := chromedp.New(ctx, chromedp.WithErrorf(logrus.Printf))
	if err != nil {
		logrus.Fatal(err)
	}

	return ins
}
