package logrushook

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisOption struct {
	Addr     string
	Password string
	DB       int
	Key      string
}

type RedisHook struct {
	appName string
	option  *RedisOption
	rClient *redis.Client
	logWash LogWashFunc
}

func NewRedisHook(appName string, option *RedisOption, logWashFunc LogWashFunc) (*RedisHook, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
		DB:       option.DB,
	})
	if _, err := rdb.Ping().Result(); err != nil {
		return nil, err
	}
	if logWashFunc == nil {
		logWashFunc = DefaultLogWashFunc
	}
	return &RedisHook{
		appName: appName,
		rClient: rdb,
		option:  option,
		logWash: logWashFunc,
	}, nil
}

func (h *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (h *RedisHook) Fire(e *logrus.Entry) (err error) {
	bs := h.logWash(h.appName, e.Time, e.Data, e.Caller, e.Level, e.Message)
	if bs == nil {
		// ignore logs
		return nil
	}

	_, err = h.rClient.RPush(h.option.Key, bs).Result()
	return
}
