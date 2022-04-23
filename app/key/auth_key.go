package key

type CtxAuthKey struct{}

type FileUploadResponse struct {
	SecureURL string
	Err       error
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
