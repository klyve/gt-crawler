package main

import (
	"bytes"
	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/Sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const GtBackendURL = "https://gt-backend-test.herokuapp.com/insertTrack"
const TempFolder = "tmp"

type Upload struct {
	Auth auth.Authenticate
}

type UploadPayload struct {
	Private bool `json:"private"`
	File string `json:"file"`
}

func (up Upload) UploadLinks(links []string, finished chan bool, config *State) {
	token, err := up.Auth.GetToken(config.FirebaseCredentials, config.CrawlerUID, config.GoogleAPIURL)
	if err != nil {
		logrus.Error("Could not get auth token", err)
		return
	}

	for i := range links {
		if err := uploadLink(links[i], token); err != nil {
			logrus.Info(err)
			logrus.Errorf("Could not upload link: %v", links[i])
			return
		}
	}

	if err := flushTemp(); err != nil {
		logrus.Error("Could not flush temp folder")
		return
	}

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

	body, boundary, err := createMultipart(mustOpen(filePath))
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, GtBackendURL, &body)
	if err != nil {
		logrus.Error(err)
		return
	}
	req.Header.Set("Content-Type", boundary)
	req.Header.Set("token", token)

	cl := &http.Client{}
	res, err := cl.Do(req)
	if err != nil {
		return
	}

	logrus.Info(res.Status)

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
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return
	}
	// anti scraper def
	req.Header.Set("Referer", GetDomain(link))

	cl := &http.Client{}
	resp, err := cl.Do(req)
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

func GetDomain(link string) (domain string) {
	urlP, _ := url.Parse(link)

	domain = urlP.Scheme + "://" +  urlP.Host
	return
}

func createMultipart(file io.Reader) (b bytes.Buffer, boundary string, err error) {
	w := multipart.NewWriter(&b)
	var fw io.Writer

	if x, ok := file.(io.Closer); ok {
		defer x.Close()
	}

	if x, ok := file.(*os.File); ok {
		fw, err = w.CreateFormFile("file", x.Name())
		if err != nil {
			return
		}
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		return
	}


	fw, err = w.CreateFormField("field")
	if err != nil {
		return
	}

	buf := []byte("false")
	_, err = io.Copy(fw, bytes.NewBuffer(buf))
	if err != nil {
		return
	}

	defer w.Close()

	boundary = w.FormDataContentType()

	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		logrus.Fatal("Not a file", err)
	}
	return r
}