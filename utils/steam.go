package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ExtractSteamID(openIDURL string) (string, error) {
	parsedURL, err := url.Parse(openIDURL)
	if err != nil {
		return "", err
	}

	steamID := parsedURL.Path[strings.LastIndex(parsedURL.Path, "/")+1:]
	return steamID, nil
}

type UserInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	ProfileURL string `json:"profile_url"`
}

func GetUserInfo(steamID string) (*UserInfo, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", os.Getenv("STEAM_API_KEY"), steamID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Steam API'ye istek başarısız: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON çözümleme hatası: %v", err)
	}

	responseData, ok := result["response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Kullanıcı bilgileri bulunamadı")
	}

	players, ok := responseData["players"].([]interface{})
	if !ok || len(players) == 0 {
		return nil, fmt.Errorf("Kullanıcı bilgileri bulunamadı")
	}

	playerData := players[0].(map[string]interface{})

	userInfo := &UserInfo{
		ID:         playerData["steamid"].(string),
		Name:       playerData["personaname"].(string),
		Avatar:     playerData["avatar"].(string),
		ProfileURL: playerData["profileurl"].(string),
	}

	return userInfo, nil
}

func GetTotalPlaytime(steamID string) (int, error) {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_played_free_games=1&include_appinfo=1", os.Getenv("STEAM_API_KEY"), steamID)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Steam API'ye istek başarısız: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("JSON çözümleme hatası: %v", err)
	}

	responseData, ok := result["response"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("Beklenmedik JSON formatı")
	}

	gamesData, ok := responseData["games"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("Oyun bilgileri bulunamadı")
	}

	totalPlaytime := 0
	for _, game := range gamesData {
		if gameInfo, ok := game.(map[string]interface{}); ok {
			if playtime, ok := gameInfo["playtime_forever"].(float64); ok {
				totalPlaytime += int(playtime)
			}
		}
	}

	return totalPlaytime, nil
}
