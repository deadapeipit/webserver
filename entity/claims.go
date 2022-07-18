package entity

import "github.com/golang-jwt/jwt"

type MyClaims struct {
	jwt.StandardClaims
	Iat int    `json:"iat"`
	Exp int    `json:"exp"`
	Uid string `json:"uid"`
}
