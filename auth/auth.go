package auth

type Authenticate interface {
	GetToken(cPath string, uid string) (token string, err error)
}