package sites

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

// crawlXContest crawls the daily score records from xcontest.org
// https://www.xcontest.org/world/en/flights/daily-score-pg/
func CrawlXContest(ins *chromedp.CDP, ctx context.Context) (sl []string, err error) {
	var nodes []*cdp.Node

	task := getLinksFromInitialStartUrl("https://www.xcontest.org/world/en/flights/daily-score-pg/", &nodes)

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
	return
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
	if err != nil {
		return
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
