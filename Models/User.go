package Models

// Chat bot User Model
type User struct {
	Uuid string `json:"uuid"`
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry string `json:"expiry"`
}
