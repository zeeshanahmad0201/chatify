package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zeeshanahmad0201/chatify/backend/models"
)

func ValidateToken(clientToken string) (claims *models.SignedDetails, err error) {
	SECRET_KEY := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(clientToken, &models.SignedDetails{}, func(t *jwt.Token) (interface{}, error) { return []byte(SECRET_KEY), nil })

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.SignedDetails)

	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

func GenerateTokens(name string, email string, userId string) (token string, refreshToken string, err error) {
	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := &models.SignedDetails{
		Email:  email,
		Name:   name,
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &models.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("error while generating token: %s", err.Error())
		return "", "", fmt.Errorf("there was an error generating token")
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("error while generating refresh token: %s", err.Error())
		return "", "", fmt.Errorf("unable to generate the refresh token")
	}

	return token, refreshToken, nil
}

func ExtractToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func ExtractTokenFromRequest(r *http.Request) string {
	authHeaders := r.Header.Get("Authorization")
	if authHeaders == "" {
		log.Println("no auth headers found")
		return ""
	}

	return ExtractToken(authHeaders)
}
