package sessions

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisSession struct {
	CookieValue string
	ValueType   valueType
	Expire      int
}

func newRedisSession(cookieValue string, valueType valueType, expire int) Session {
	return &redisSession{
		CookieValue: cookieValue,
		ValueType:   valueType,
		Expire:      expire,
	}
}
func (r *redisSession) stringer() string {
	return prefixValue + r.CookieValue
}
func (r *redisSession) Get(keys ...string) interface{} {

	if len(keys) == 0 && r.ValueType == ValueMap {
		panic("key can't nil")
	}

	if r.ValueType == ValueMap {
		result, err := rdb.HGet(context.TODO(), r.stringer(), keys[0]).Result()
		if err != nil && err != redis.Nil {
			panic(fmt.Sprintf("hget value error:%s", err))
		}
		return result
	}
	result, err := rdb.Get(context.TODO(), r.stringer()).Result()
	if err != nil && err != redis.Nil {
		panic(fmt.Sprintf("get value error:%s", err))
	}
	return result
}

func (r *redisSession) Set(value interface{}, keys ...string) {
	if r.ValueType == ValueMap && len(keys) == 0 {
		panic("The type is map. Please pass a key")
	} else if r.ValueType == ValueString && len(keys) > 0 {
		panic("The type is string. Please don't pass the key")
	}

	if r.ValueType == ValueMap {
		if err := rdb.HSet(context.TODO(), r.stringer(), keys[0], value).Err(); err != nil {
			panic(fmt.Sprintf("hset value error:%s", err.Error()))
		}
	} else {
		if err := rdb.Set(context.TODO(), r.stringer(), value, time.Duration(0)).Err(); err != nil {
			panic(fmt.Sprintf("set value error:%s", err.Error()))
		}
	}
	rdb.Expire(context.TODO(), r.stringer(), time.Second*time.Duration(r.Expire))
}

func (r *redisSession) Delete(keys ...string) {
	rdb.Del(context.TODO(), keys...)
}
