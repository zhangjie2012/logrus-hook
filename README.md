# logrus-hook

Go [logrus](https://github.com/sirupsen/logrus) hooks.

```bash
go get github.com/zhangjie2012/logrus-hook
```

## Redis LIST

log write(`RPUSH`) to redis LIST.

```go
option := RedisOption{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    Key:      "logrusredis.hook",
}
hook, _ := NewRedisHook("application_name", &option, nil)
logrus.AddHook(hook)
```

## Customize

If you want customize inserted redis bs, you can customize a `LogWashFunc`, checkout `redishook_test.go` file.
