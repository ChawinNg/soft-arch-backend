package model

import "github.com/golang-jwt/jwt"

type SessionClaims struct {
	jwt.StandardClaims
	Session Sessions `json:"session"`
}
