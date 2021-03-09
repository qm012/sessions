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
import (
	"github.com/gin-gonic/gin"
	"github.com/qm012/sessions"
)

var (
	UsernameStr = "username"
	PasswordStr = "password"
	IsLoginStr  = "isLogin"
)
r := gin.Default()
sessions.SetCookie(sessions.ValueMap, "session_id", age, "/", "127.0.0.1", false, false)
sessions.SetCookie(sessions.ValueString, "password", age, "/", "127.0.0.1", false, false)
r.Use(sessions.Sessions(sessions.ChooseSessionStore(sessions.Redis, rdbClient)))

r.GET("/login",func(c *gin.Context){
    sessions.GetSession(ctx, "password").Set(user.Password)
    sessions.GetSession(ctx, "session_id").Set(true, IsLoginStr)
    sessions.GetSession(ctx, "session_id").Set(user.Username, UsernameStr)
    sessions.GetSession(ctx, "session_id").Set(user.Password, PasswordStr)
})
r.GET("/home",func(c *gin.Context){
    username := sessions.GetSession(ctx, "session_id").Get(UsernameStr)
	password := sessions.GetSession(ctx, "password").Get()
	login := sessions.GetSession(ctx, "session_id").Get(IsLoginStr)
})
r.Run()
```

### 内存版本

```go
r := gin.Default()
sessions.SetCookie(sessions.ValueMap, "session_id", age, "/", "127.0.0.1", false, false)
sessions.SetCookie(sessions.ValueString, "password", age, "/", "127.0.0.1", false, false)

r.Use(sessions.Sessions(sessions.ChooseSessionStore(sessions.Memory)))

r.GET("/login",func(c *gin.Context){
    sessions.GetSession(ctx, "password").Set(user.Password)
    sessions.GetSession(ctx, "session_id").Set(true, IsLoginStr)
    sessions.GetSession(ctx, "session_id").Set(user.Username, UsernameStr)
    sessions.GetSession(ctx, "session_id").Set(user.Password, PasswordStr)
})
r.GET("/home",func(c *gin.Context){
    username := sessions.GetSession(ctx, "session_id").Get(UsernameStr)
	password := sessions.GetSession(ctx, "password").Get()
	login := sessions.GetSession(ctx, "session_id").Get(IsLoginStr)
})
r.Run()

```