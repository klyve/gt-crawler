package auth

import (
	"bytes"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
)


type FAuth struct{}

type authResponse struct {
	Kind         string
	IDToken      string
	RefreshToken string
	ExpiresIn    string
}

func (auth FAuth) GetToken(cPath string, uid string, urlGoogleAPI string) (token string, err error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return
	}

	token, err = client.CustomToken(ctx, uid)
	if err != nil {
		return
	}

	token, err = client.CustomToken(context.Background(), uid)
	if err != nil {
		logrus.Fatalf("error setting custom token: %v\n", err)
	}

	var jsonStr = []byte(`{
	"token": "` + token + `",
	"returnSecureToken": true
	}`)
	res, err := http.Post(urlGoogleAPI, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Printf("%+v\n", string(jsonStr))
		logrus.Fatalf("error retrieving id token: %v\n", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Fatalf("error reading body of id token retrieve req: %v\n", err)
	}

	var resParsed authResponse
	err = json.Unmarshal(resBody, &resParsed)
	if err != nil {
		logrus.Fatalf("error parsing json of id token request: %v\n", err)
	}

	token = resParsed.IDToken
	logrus.Infof("Firebase token: %v", resParsed.IDToken)

	return
}
