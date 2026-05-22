package token

import (
	"errors"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/middleware/my_jwt"
	"ginskeleton/app/model"
	"ginskeleton/app/service/users/token_cache_redis"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateUserFactory() *userToken {
	return &userToken{
		userJwt: my_jwt.CreateMyJWT(variable.ConfigYml.GetString("Token.JwtTokenSignKey")),
	}
}

type userToken struct {
	userJwt *my_jwt.JwtSign
}

func (u *userToken) GenerateToken(userid int64, username string, phone string, expireAt int64) (tokens string, err error) {

	customClaims := my_jwt.CustomClaims{
		UserId: userid,
		Name:   username,
		Phone:  phone,

		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,
			ExpiresAt: time.Now().Unix() + expireAt,
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userId := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return model.CreateUserFactory("").OauthLoginToken(userId, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}

func (u *userToken) TokenIsMeetRefreshCondition(token string) bool {

	customClaims, code := u.isNotExpired(token, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshAllowSec"))
	switch code {
	case consts.JwtTokenOK, consts.JwtTokenExpired:

		if model.CreateUserFactory("").OauthRefreshConditionCheck(customClaims.UserId, token) {
			return true
		}
	}
	return false
}

func (u *userToken) RefreshToken(oldToken, clientIp string) (newToken string, res bool) {
	var err error

	if newToken, err = u.userJwt.RefreshToken(oldToken, variable.ConfigYml.GetInt64("Token.JwtTokenRefreshExpireAt")); err == nil {
		if customClaims, err := u.userJwt.ParseToken(newToken); err == nil {
			userId := customClaims.UserId
			expiresAt := customClaims.ExpiresAt
			if model.CreateUserFactory("").OauthRefreshToken(userId, expiresAt, oldToken, newToken, clientIp) {
				return newToken, true
			}
		}
	}

	return "", false
}

func (u *userToken) isNotExpired(token string, expireAtSec int64) (*my_jwt.CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {

		if time.Now().Unix()-(customClaims.ExpiresAt+expireAtSec) < 0 {

			return customClaims, consts.JwtTokenOK
		} else {

			return customClaims, consts.JwtTokenExpired
		}
	} else {

		return nil, consts.JwtTokenInvalid
	}
}

func (u *userToken) IsEffective(token string) bool {
	customClaims, code := u.isNotExpired(token, 0)
	if consts.JwtTokenOK == code {

		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			tokenRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(customClaims.UserId)
			if tokenRedisFact != nil {
				defer tokenRedisFact.ReleaseRedisConn()
				if tokenRedisFact.TokenCacheIsExists(token) {
					return true
				}
			}
		}

		if model.CreateUserFactory("").OauthCheckTokenIsOk(customClaims.UserId, token) {
			return true
		}
	}
	return false
}

func (u *userToken) ParseToken(tokenStr string) (CustomClaims my_jwt.CustomClaims, err error) {
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil {
		return *customClaims, nil
	} else {
		return my_jwt.CustomClaims{}, errors.New(my_errors.ErrorsParseTokenFail)
	}
}

func (u *userToken) DestroyToken() {

}
