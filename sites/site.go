package sites

import (
	"context"
	"github.com/chromedp/chromedp"
)

type ChromeSite interface {
	Crawl(c *chromedp.CDP, ctx context.Context) (sl []string, err error)
}