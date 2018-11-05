package main

import (
	authenticate "github.com/GlidingTracks/gt-crawler/auth"
	"github.com/Sirupsen/logrus"
)

func uploadLinks(links []string, finished chan bool, cPath string, uid string) {

	auth := authenticate.FAuth{}

	token, err := auth.GetFToken(cPath, uid)
	if err != nil {
		logrus.Error("Could not get firebase token")
		return
	}

	// Upload finished
	finished <- true
}
