# sessions
go gin版本的session存储，支持存取多个cookie，内存版和redis版本


## 使用

下载

```bash
$ go get github.com/qm012/sessions
```

导入

```go
import "github.com/qm012/sessions"
```

## 示例

示例项目 [session-demo](https://github.com/qm012/sessions-demo)

### redis 版本
```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/qm012/sessions"
	"log"
	"net/http"
)

var (
	UsernameStr = "username"
	PasswordStr = "password"
	IsLoginStr  = "isLogin"
	rdbClient   *redis.Client
)

func init() {
	rdbClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
		PoolSize: 1000,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := rdbClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	log.Print("redis link success")
}

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func main() {

	r := gin.Default()
	var age = 60
	sessions.SetCookie(sessions.ValueMap, "session_id", age, "/", "127.0.0.1", false, false)
	sessions.SetCookie(sessions.ValueString, "password", age, "/", "127.0.0.1", false, false)
	r.Use(sessions.Sessions(sessions.ChooseSessionStore(sessions.Redis, rdbClient)))

	r.GET("/login", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		}
		sessions.GetSession(ctx, "password").Set(user.Password)
		sessions.GetSession(ctx, "session_id").Set(true, IsLoginStr)
		sessions.GetSession(ctx, "session_id").Set(user.Username, UsernameStr)
		sessions.GetSession(ctx, "session_id").Set(user.Password, PasswordStr)
	})
	r.GET("/home", func(ctx *gin.Context) {
		username := sessions.GetSession(ctx, "session_id").Get(UsernameStr)
		ctx.HTML(http.StatusOK, "home.html", gin.H{"username": username})
	})
	r.Run()
}

```

### 内存版本

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/qm012/sessions"
	"net/http"
)

var (
	UsernameStr = "username"
	PasswordStr = "password"
	IsLoginStr  = "isLogin"
	rdbClient   *redis.Client
)

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func main() {

	r := gin.Default()
	var age = 49
	sessions.SetCookie(sessions.ValueMap, "session_id", age, "/", "127.0.0.1", false, false)
	sessions.SetCookie(sessions.ValueString, "password", age, "/", "127.0.0.1", false, false)

	r.Use(sessions.Sessions(sessions.ChooseSessionStore(sessions.Memory)))

	r.GET("/login", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		}
		sessions.GetSession(ctx, "password").Set(user.Password)
		sessions.GetSession(ctx, "session_id").Set(true, IsLoginStr)
		sessions.GetSession(ctx, "session_id").Set(user.Username, UsernameStr)
		sessions.GetSession(ctx, "session_id").Set(user.Password, PasswordStr)
	})
	r.GET("/home", func(ctx *gin.Context) {
		username := sessions.GetSession(ctx, "session_id").Get(UsernameStr)
		ctx.HTML(http.StatusOK, "home.html", gin.H{"username": username})
	})
	r.Run()
}

```