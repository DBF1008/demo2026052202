package redis_factory

import (
	"ginskeleton/app/core/event_manage"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/yml_config"
	"ginskeleton/app/utils/yml_config/ymlconfig_interf"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var redisPool *redis.Pool
var configYml ymlconfig_interf.YmlConfigInterf

func init() {
	configYml = yml_config.CreateYamlFactory()
	redisPool = initRedisClientPool()
}
func initRedisClientPool() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     configYml.GetInt("Redis.MaxIdle"),
		MaxActive:   configYml.GetInt("Redis.MaxActive"),
		IdleTimeout: configYml.GetDuration("Redis.IdleTimeout") * time.Second,
		Dial: func() (redis.Conn, error) {

			conn, err := redis.Dial("tcp", configYml.GetString("Redis.Host")+":"+configYml.GetString("Redis.Port"))
			if err != nil {
				variable.ZapLog.Error(my_errors.ErrorsRedisInitConnFail + err.Error())
				return nil, err
			}
			auth := configYml.GetString("Redis.Auth")
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					_ = conn.Close()
					variable.ZapLog.Error(my_errors.ErrorsRedisAuthFail + err.Error())
				}
			}
			_, _ = conn.Do("select", configYml.GetInt("Redis.IndexDb"))
			return conn, err
		},
	}

	eventManageFactory := event_manage.CreateEventManageFactory()
	if _, exists := eventManageFactory.Get(variable.EventDestroyPrefix + "Redis"); exists == false {
		eventManageFactory.Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
			_ = redisPool.Close()
		})
	}
	return redisPool
}

func GetOneRedisClient() *RedisClient {
	maxRetryTimes := configYml.GetInt("Redis.ConnFailRetryTimes")
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = redisPool.Get()

		if _, replyErr := oneConn.Do("time"); replyErr != nil {

			initRedisClientPool()
			oneConn = redisPool.Get()
		}

		if err := oneConn.Err(); err != nil {

			if i == maxRetryTimes {
				variable.ZapLog.Error(my_errors.ErrorsRedisGetConnFail, zap.Error(oneConn.Err()))
				return nil
			}

			time.Sleep(time.Second * configYml.GetDuration("Redis.ReConnectInterval"))
		} else {
			break
		}
	}
	return &RedisClient{oneConn}
}

type RedisClient struct {
	client redis.Conn
}

func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

func (r *RedisClient) ReleaseOneRedisClient() {
	_ = r.client.Close()
}

func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func (r *RedisClient) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

func (r *RedisClient) StringMap(reply interface{}, err error) (map[string]string, error) {
	return redis.StringMap(reply, err)
}

func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

func (r *RedisClient) IntMap(reply interface{}, err error) (map[string]int, error) {
	return redis.IntMap(reply, err)
}

func (r *RedisClient) Int64Map(reply interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(reply, err)
}

func (r *RedisClient) Int64s(reply interface{}, err error) ([]int64, error) {
	return redis.Int64s(reply, err)
}

func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

