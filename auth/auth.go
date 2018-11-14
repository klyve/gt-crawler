package auth

type Authenticate interface {
	GetToken(cPath string, uid string, url string) (token string, err error)
}