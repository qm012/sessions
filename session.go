package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"sync"
)

// cookiesMap store manage cookie data
var (
	cookies []*cookieType
	once    sync.Once
)

type cookieType struct {
	cookie    *http.Cookie
	valueType valueType
}

const (
	Redis  storeType = "redis"
	Memory storeType = "memory"

	ValueMap valueType = iota
	ValueString
)

type (
	// data store type
	storeType string
	// cookie value type
	valueType int
)

// Session single cookie value object
type Session interface {
	Get(keys ...string) interface{}
	Set(value interface{}, keys ...string)
	Delete(keys ...string)
}

// SessionMgr object manage
type SessionMgr interface {
	GetSession(cookValue string) Session
	CreateSession(cookValue string, valueType valueType, expire int) Session
}

func ChooseSessionStore(st storeType, clients ...*redis.Client) (mgr SessionMgr) {
	if st == "" {
		st = Memory
	}
	if st == Redis && len(clients) == 0 {
		panic("need clients")
	}
	switch st {
	case Redis:
		mgr = newRedisMgr(clients[0])
	case Memory:
		mgr = newMemoryMgr()
	default:
		mgr = newMemoryMgr()
		return
	}
	return
}

// SetCookie setting cookie
func SetCookie(valueType valueType, name string, maxAge int, path, domain string, secure, httpOnly bool) {
	once.Do(func() {
		cookies = make([]*cookieType, 0, 10)
	})
	if name == "" {
		panic("cookie name can't nil")
	}

	for _, v := range cookies {
		if v.cookie.Name == name {
			panic("Please check if you have the same cookie name")
		}
	}
	cookies = append(cookies, &cookieType{
		valueType: valueType,
		cookie: &http.Cookie{
			Name:     name,
			MaxAge:   maxAge,
			Path:     path,
			Domain:   domain,
			Secure:   secure,
			HttpOnly: httpOnly,
		},
	})
}

// Sessions middleware of cookies
func Sessions(mgr SessionMgr) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if len(cookies) == 0 {
			return
		}

		// request cookies
		rCookies := ctx.Request.Cookies()

		for _, cookie := range cookies {
			var session Session
			value := uuid.NewV4().String()
			// search by request cookies
			for _, rCookie := range rCookies {
				if cookie.cookie.Name == rCookie.Name {
					session = mgr.GetSession(rCookie.Value)
					value = rCookie.Value
					break
				}
			}
			// not found

			if session == nil {
				session = mgr.CreateSession(value, cookie.valueType, cookie.cookie.MaxAge)
			}
			// set gin key value is session
			ctx.Set(cookie.cookie.Name, session)

			//The cookie must be set before the handler returns
			ctx.SetCookie(cookie.cookie.Name, url.QueryEscape(value), cookie.cookie.MaxAge,
				cookie.cookie.Path, cookie.cookie.Domain, cookie.cookie.Secure, cookie.cookie.HttpOnly)
		}

		// execute net gin.HandlerFuc
		ctx.Next()
	}
}

// GetSession Session Object
// If you use a custom plaintext cookie, use c.Cookie(name) or http.Request.Cookie(name)
// @param ctx  gin.context
// @param name is cookie name
func GetSession(ctx *gin.Context, name string) Session {
	return ctx.MustGet(name).(Session)
}
