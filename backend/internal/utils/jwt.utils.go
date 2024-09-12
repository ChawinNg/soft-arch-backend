package utils

import (
	"backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJwtToken(session model.Sessions, maxAge time.Duration, jwtSecret string) (string, error) {
	issuedTime := time.Now().Unix()
	expiresTime := time.Now().Add(maxAge).Unix()

	customClaim := model.SessionClaims{
		Session: session,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedTime,
			ExpiresAt: expiresTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(token string, jwtSecret string) (*model.SessionClaims, error) {
	claim, err := jwt.ParseWithClaims(token, &model.SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claim.Claims.(*model.SessionClaims), nil
}
