package main

func uploadLinks(links []string, finished chan bool) {


	// Upload finished
	finished <- true
}
