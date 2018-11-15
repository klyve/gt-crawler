package sites

import (
	"context"
)

type ChromeSite interface {
	Crawl(ctx context.Context) (sl []string, err error)
}
