package key

type CtxAuthKey struct{}

type FileUploadResponse struct {
	SecureURL string
	Err       error
}

type Token struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
