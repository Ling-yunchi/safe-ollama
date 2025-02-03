package model

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtPayload struct {
	UserId  uint   `json:"userId"`
	Role    string `json:"role"`
	Expires int64  `json:"expires"`
}

func (p *JwtPayload) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"userId":  p.UserId,
		"role":    p.Role,
		"expires": p.Expires,
	}
}

func (p *JwtPayload) FromMapClaims(claims jwt.MapClaims) {
	p.UserId = uint(claims["userId"].(float64))
	p.Role = claims["role"].(string)
	p.Expires = int64(claims["expires"].(float64))
}

type BaseResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
