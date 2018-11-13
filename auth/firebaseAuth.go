package auth

import (
	"firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)


type FAuth struct{}

func (auth FAuth) GetToken(cPath string, uid string) (token string, err error) {
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

	return
}
