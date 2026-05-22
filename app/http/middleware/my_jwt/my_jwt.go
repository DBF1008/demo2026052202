package my_jwt

import (
	"errors"
	"ginskeleton/app/global/my_errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateMyJWT(signKey string) *JwtSign {
	if len(signKey) <= 0 {
		signKey = "ginskeleton"
	}
	return &JwtSign{
		[]byte(signKey),
	}
}

type JwtSign struct {
	SigningKey []byte
}

func (j *JwtSign) CreateToken(claims CustomClaims) (string, error) {

	tokenPartA := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenPartA.SignedString(j.SigningKey)
}

func (j *JwtSign) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if token == nil {
		return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(my_errors.ErrorsTokenMalFormed)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(my_errors.ErrorsTokenNotActiveYet)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {

				token.Valid = true
				goto labelHere
			} else {
				return nil, errors.New(my_errors.ErrorsTokenInvalid)
			}
		}
	}
labelHere:
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
}

func (j *JwtSign) RefreshToken(tokenString string, extraAddSeconds int64) (string, error) {

	if CustomClaims, err := j.ParseToken(tokenString); err == nil {
		CustomClaims.ExpiresAt = time.Now().Unix() + extraAddSeconds
		return j.CreateToken(*CustomClaims)
	} else {
		return "", err
	}
}
