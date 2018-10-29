package main

type Crawler interface {
	Crawl(url string)
}

func main() {
	c := &Chrome{}

	c.Crawl()
}
