package sessions

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb *redis.Client

type redisMgr struct {
}

func newRedisMgr(client *redis.Client) SessionMgr {
	rdb = client
	return &redisMgr{}
}

const (
	prefixData  = "cookie:data:"
	prefixValue = "cookie:value:"
)

func (r *redisMgr) GetSession(cookValue string) (session Session) {
	buf, err := rdb.Get(context.TODO(), prefixData+cookValue).Bytes()
	if err != nil && err != redis.Nil {
		panic(fmt.Sprintf("get data error: %s", err.Error()))
	}
	if len(buf) == 0 {
		return nil
	}

	var redisSession redisSession
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err = dec.Decode(&redisSession)

	if err != nil {
		panic(fmt.Sprintf("err:%s", err.Error()))
	}
	return &redisSession
}

func (r *redisMgr) CreateSession(cookValue string, valueType valueType, expire int) (session Session) {

	session = newRedisSession(cookValue, valueType, expire)

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(session)

	if err != nil {
		panic(fmt.Sprintf("Encode session err:%s", err))
	}

	_, err = rdb.Set(context.TODO(), prefixData+cookValue, buf.Bytes(), time.Duration(0)).Result()
	rdb.Expire(context.TODO(), prefixData+cookValue, time.Second*time.Duration(expire))
	if err != nil {
		panic(fmt.Sprintf("set []byte session err:%s", err))
	}
	return
}
