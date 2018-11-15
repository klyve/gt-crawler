package main

import (
	"bytes"
	"context"
	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/Sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GtBackendURL holds the server url which to upload to.
const GtBackendURL = "https://gt-backend-test.herokuapp.com/insertTrack"

// TempFolder name of folder to temporary store files.
const TempFolder = "tmp"

// Upload struct holds authentication object. Determines which method to use to get tokens.
type Upload struct {
	Auth auth.Authenticate
}

// UploadLinks will download files and upload those files to the server.
func (up Upload) UploadLinks(ctx context.Context, links []string, config *State) (finished bool) {
	token, err := up.Auth.GetToken(ctx, config.FirebaseCredentials, config.CrawlerUID, config.GoogleAPIURL)
	if err != nil {
		logrus.Error("Could not get auth token", err)
		return
	}

	wg := sync.WaitGroup{}
	for i := range links {
		wg.Add(1)
		go func(j int) {
			if err := uploadLink(links[j], token); err != nil {
				logrus.Errorf("Could not upload link: %v | error: %v", links[j], err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if err := flushTemp(); err != nil {
		logrus.Error("Could not flush temp folder")
		return
	}

	finished = true
	return
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

	if res.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		body := buf.String()
		logrus.Errorf("Error in uploading: file: %v, status: %v, message: %v", filePath, res.Status, body)
	}

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

// GetFileName will strip characters which should be escaped, does not strip special characters per say.
func GetFileName(link string) (fileName string) {
	parts := strings.Split(link, "www.")

	name := parts[len(parts)-1]

	fileName = strings.Replace(name, "/", "_", -1)

	return
}

// GetDomain will return the base address from a url. e.g., http://www.test.com/api/info/add/ => http://www.test.com.
func GetDomain(link string) (domain string) {
	urlP, _ := url.Parse(link)

	domain = urlP.Scheme + "://" + urlP.Host
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
