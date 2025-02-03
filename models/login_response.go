package models

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
