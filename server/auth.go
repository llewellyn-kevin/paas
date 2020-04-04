package main

import(
  "fmt"
  "log"
  "net/http"

  // "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

const(
  AUTHENTICATED = "auth"
  ANONYMOUS = "anom"
)

func Authorize() gin.HandlerFunc {
  return func(c *gin.Context) {
    cookie, err := c.Cookie("auth_jwt")

    if err != nil {
      c.Set("authorization", ANONYMOUS)
      c.Set("user", "")
    } else {
      c.Set("authorization", AUTHENTICATED)
      c.Set("user", cookie)
    }
  }
}

func Login(c *gin.Context, user string) {
  log.Println("setting cookie foe user")
  c.SetCookie("auth_jwt", user, 3600, "/", "127.0.0.1", http.SameSiteNoneMode, false, false)
}

func IsValidAuth(c *gin.Context, store redis.Conn, user, pass string) bool {
  db_pass, err := redis.String(store.Do("HGET", fmt.Sprintf("user:%s", user), "password"))
  if err != nil {
    return false
  } else {
    return db_pass == pass
  }
}
