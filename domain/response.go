package domain

type ErrorResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
