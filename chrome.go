package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

type Chrome struct{}

func (ch Chrome) Crawl() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	c := createInstance(ctx)

	sl := crawlXContest(c, ctx)

	//Upload source links
	for i := range sl {
		logrus.Info(sl[i])
	}
}

// crawlXContest crawls the daily score records from xcontest.org
// https://www.xcontest.org/world/en/flights/daily-score-pg/
func crawlXContest(ins *chromedp.CDP, ctx context.Context) (sl []string) {
	var nodes []*cdp.Node

	task := getLinksFromInitialStartUrl("https://www.xcontest.org/world/en/flights/daily-score-pg/", &nodes)

	err := ins.Run(ctx, task)
	if err != nil {
		logrus.Error(err)
	}

	err = ins.Shutdown(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ins.Wait()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Shutting down initial instance")

	visitQueue := filterRealFlightLinks(nodes)

	sl = visitDetailsPagesAndExtract(visitQueue, ctx)
	return
}

func visitDetailsPagesAndExtract(urls []string, ctx context.Context) (sl []string) {
	ins := createInstance(ctx)

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

	err := ins.Shutdown(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ins.Wait()
	if err != nil {
		logrus.Fatal(err)
	}

	return
}

func getLinksFromInitialStartUrl(url string, nodes *[]*cdp.Node) chromedp.Tasks {
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

func createInstance(ctx context.Context) (ins *chromedp.CDP) {
	ins, err := chromedp.New(ctx, chromedp.WithErrorf(logrus.Printf))
	if err != nil {
		logrus.Fatal(err)
	}

	return ins
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
