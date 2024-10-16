package domain

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	SteamID string `json:"steam_id"`
}
