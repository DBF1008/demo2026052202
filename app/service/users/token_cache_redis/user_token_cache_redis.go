package token_cache_redis

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/md5_encrypt"
	"ginskeleton/app/utils/redis_factory"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func CreateUsersTokenCacheFactory(userId int64) *userTokenCacheRedis {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userTokenCacheRedis{redisClient: redCli, userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10)}
}

type userTokenCacheRedis struct {
	redisClient  *redis_factory.RedisClient
	userTokenKey string
}

func (u *userTokenCacheRedis) SetTokenCache(tokenExpire int64, token string) bool {

	if _, err := u.redisClient.Int(u.redisClient.Execute("zAdd", u.userTokenKey, tokenExpire, md5_encrypt.MD5(token))); err == nil {
		return true
	} else {
		variable.ZapLog.Error("缓存用户token到redis出错", zap.Error(err))
	}
	return false
}

func (u *userTokenCacheRedis) DelOverMaxOnlineCache() bool {

	_, _ = u.redisClient.Execute("zRemRangeByScore", u.userTokenKey, 0, time.Now().Unix()-1)

	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	alreadyCacheNum, err := u.redisClient.Int(u.redisClient.Execute("zCard", u.userTokenKey))
	if err == nil && alreadyCacheNum > onlineUsers {

		if alreadyCacheNum, err = u.redisClient.Int(u.redisClient.Execute("zRemRangeByRank", u.userTokenKey, 0, alreadyCacheNum-onlineUsers-1)); err == nil {
			return true
		} else {
			variable.ZapLog.Error("删除超过系统允许之外的token出错：", zap.Error(err))
		}
	}
	return false
}

func (u *userTokenCacheRedis) TokenCacheIsExists(token string) (exists bool) {
	token = md5_encrypt.MD5(token)
	curTimestamp := time.Now().Unix()
	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if strSlice, err := u.redisClient.Strings(u.redisClient.Execute("zRevRange", u.userTokenKey, 0, onlineUsers-1)); err == nil {
		for _, val := range strSlice {
			if score, err := u.redisClient.Int64(u.redisClient.Execute("zScore", u.userTokenKey, token)); err == nil {
				if score > curTimestamp {
					if strings.Compare(val, token) == 0 {
						exists = true
						break
					}
				}
			}
		}
	} else {
		variable.ZapLog.Error("获取用户在redis缓存的 token 值出错：", zap.Error(err))
	}
	return
}

func (u *userTokenCacheRedis) SetUserTokenExpire(ts int64) bool {
	if _, err := u.redisClient.Execute("expireAt", u.userTokenKey, ts); err == nil {
		return true
	}
	return false
}

func (u *userTokenCacheRedis) ClearUserToken() bool {
	if _, err := u.redisClient.Execute("del", u.userTokenKey); err == nil {
		return true
	}
	return false
}

func (u *userTokenCacheRedis) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
