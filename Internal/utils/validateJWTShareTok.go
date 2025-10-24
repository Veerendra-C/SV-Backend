package utils

import (
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateTokens(tokenStr string) (*modals.FileShareClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &modals.FileShareClaims{}, func(t *jwt.Token) (interface{}, error) {
		return SecretShareKey,nil
	})
	if err != nil{
		return nil, err
	}

	if claims, ok := token.Claims.(*modals.FileShareClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}