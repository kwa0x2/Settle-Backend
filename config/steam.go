package config

import (
	"fmt"
	"net/url"
	"os"
)

func GetSteamLoginURL() string {
	return fmt.Sprintf(
		"https://steamcommunity.com/openid/login?openid.mode=checkid_setup&openid.ns=http://specs.openid.net/auth/2.0&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select&openid.identity=http://specs.openid.net/auth/2.0/identifier_select&openid.return_to=%s",
		url.QueryEscape(os.Getenv("STEAM_REDIRECT_URL")),
	)
}
