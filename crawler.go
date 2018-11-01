package main

type Crawler interface {
	Crawl()
}

func main() {
	c := &Chrome{}

	c.Crawl()
}
