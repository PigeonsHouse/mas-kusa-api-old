package routers

import "github.com/golang-jwt/jwt/v4"

type jwtCustomClaims struct {
	Instance string `json:"instance"`
	Token    string `json:"token"`
	jwt.StandardClaims
}

type SimpleJson struct {
	Message string `json:"message"`
}

type ImagePath struct {
	Path    string `json:"path"`
	Refresh bool   `json:"refresh"`
}

type Register struct {
	Instance string `json:"instance"`
	Token    string `json:"token"`
}

type Token struct {
	Token string `json:"token"`
}
