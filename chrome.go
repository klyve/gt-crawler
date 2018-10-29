package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strings"
)

type Chrome struct{}

func (ch Chrome) Crawl() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := chromedp.New(ctx,
		chromedp.WithErrorf(logrus.Printf),
		/*chromedp.WithRunnerOptions(
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true))*/)
	if err != nil {
		logrus.Fatal(err)
	}

	crawlXContest(c, ctx)
}

// crawlXContest crawls the daily score records from xcontest.org
// https://www.xcontest.org/world/en/flights/daily-score-pg/
func crawlXContest(ins *chromedp.CDP, ctx context.Context) {
	// #flight-1539343 > td:nth-child(11) > div:nth-child(1) > a:nth-child(1)
	// #flight-1539310 > td:nth-child(11) > div:nth-child(1) > a:nth-child(1)

	// Table selector: XClist

	var nodes []*cdp.Node

	task := getLinksToFollow("https://www.xcontest.org/world/en/flights/daily-score-pg/", &nodes)

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

	visitQueue := filterRealFlightLinks(nodes)
	for i := range visitQueue {
		logrus.Info(visitQueue[i])
	}
}

func getLinksToFollow(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#flights`),
		chromedp.WaitReady(`//a`, chromedp.BySearch),
		chromedp.Nodes(`//a`, nodes, chromedp.BySearch),
	}
}

func filterRealFlightLinks(found []*cdp.Node) (visitQueue []string){
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