package auth

import "context"

type Authenticate interface {
	GetToken(ctx context.Context, cPath string, uid string, url string) (token string, err error)
}
