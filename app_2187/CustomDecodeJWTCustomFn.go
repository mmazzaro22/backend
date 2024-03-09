package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// DecodeJWTCustomFn is a custom function.
func DecodeJWTCustomFn(jwttoken string) (success bool, userID int, err error) {

	token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret_key"), nil
	})
	if err != nil || token.Valid != true {
		err = nil
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if foundId, found := claims["user_id"]; found {
			if idInt, ok := foundId.(float64); ok {
				userID = int(idInt)
				success = true
				return
			} else {
				return
			}
		} else {
			return
		}
	}

	return

}
