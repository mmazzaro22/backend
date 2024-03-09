package main

import (
	"github.com/dgrijalva/jwt-go"
)

// EncodeJWTCustomFn is a custom function.
func EncodeJWTCustomFn(userID int) (token string, err error) {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	tok, err := jwtToken.SignedString([]byte("secret_key"))
	if err != nil {
		return
	}

	return tok, nil

}
