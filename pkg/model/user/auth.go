package user

type GoogleOAuth2Code struct {
	Code string `json:"code" binding:"required"`
}
