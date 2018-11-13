package main

import (
	"bytes"
	"encoding/json"
	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const GtBackendURL = "https://gt-backend-test.herokuapp.com/"
const TempFolder = "tmp"

type Upload struct {
	Auth auth.Authenticate
}

type UploadPayload struct {
	Token string `json:"token"`
	Private bool `json:"private"`
	File string `json:"file"`
}

func (up Upload) UploadLinks(links []string, finished chan bool, cPath string, uid string) {
	token, err := up.Auth.GetToken(cPath, uid)
	if err != nil {
		logrus.Error("Could not get auth token")
		return
	}

	for i := range links {
		if err := uploadLink(links[i], token); err != nil {
			logrus.Errorf("Could not upload link: %v", links[i])
			return
		}
	}

	/*if err := flushTemp(); err != nil {
		logrus.Error("Could not flush temp folder")
		return
	}*/

	// Upload finished
	if finished != nil {
		finished <- true
	}
}

func flushTemp() (err error) {
	err = os.RemoveAll(TempFolder)
	return
}

func uploadLink(link string, token string) (err error) {
	filePath, err := downloadFile(link)
	if err != nil {
		logrus.Error(err)
		return
	}

	body, err := createBody(filePath, token)
	if err != nil {
		logrus.Error(err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, GtBackendURL, bytes.NewBuffer(body))
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info(req)

	return
}

func downloadFile(link string) (path string, err error) {
	fileName := GetFileName(link)
	path = filepath.Join(TempFolder, fileName)
	// Create the file

	os.MkdirAll(TempFolder, os.ModePerm)

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}

func GetFileName(link string) (fileName string) {
	parts := strings.Split(link, "www.")

	name := parts[len(parts)-1]

	fileName = strings.Replace(name, "/", "_", -1)

	return
}

func createBody(filePath string, token string) (body []byte, err error) {
	upl := &UploadPayload{
		Token: token,
		Private: false,
		File: filePath,
	}

	body, err = json.Marshal(upl)
	return
}