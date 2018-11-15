package main

// State keeps track of config vars associated with the crawler
type State struct {
	DisallowedDomains   []string
	CrawlerUID          string
	FirebaseCredentials string
	GoogleAPIURL        string
}
