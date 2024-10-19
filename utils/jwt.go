package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"time"
)

type Claims struct {
	User *domain.User `json:"user"`
	jwt.RegisteredClaims
}

func CreateAccessToken(user *domain.User, secret string, exp int) (string, error) {
	claims := Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(exp) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user *domain.User, secret string, exp int) (string, error) {
	claims := Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(exp) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (*domain.User, error) {
	parsedToken, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("Invalid token claims or token is not valid")
	}

	userData, ok := claims["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("User data not found in claims")
	}

	user := &domain.User{
		ID:         userData["ID"].(string),
		Name:       userData["Name"].(string),
		Avatar:     userData["Avatar"].(string),
		ProfileURL: userData["ProfileURL"].(string),
		Role:       types.UserRole(userData["Role"].(string)),
	}
	return user, nil
}
