package main

import(
  "fmt"
  "io/ioutil"
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

const(
  ServiceName = "PaaS"
  TokenLifetime = 360

  Authenticated = "auth"
  Anonymous = "anoma"
)

func GetSecret(method string) ([]byte, error) {
  return ioutil.ReadFile(fmt.Sprintf("jwt_secret.%s", method))
}

type AuthClaims struct {
  Username string `json:"usr"`
  jwt.StandardClaims
}

func GetAuthClaims(user string) AuthClaims {
  now := time.Now().Unix()
  return AuthClaims{
    user,
    jwt.StandardClaims{
      Issuer: ServiceName,
      IssuedAt: now,
      ExpiresAt: now + TokenLifetime,
    },
  }
}

func Authorize() gin.HandlerFunc {
  return func(c *gin.Context) {
    cookie, err := c.Cookie("auth_jwt")

    if err != nil {
      c.Set("authorization", Anonymous)
      c.Set("user", "")
    } else {
      token, err := jwt.ParseWithClaims(cookie, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
        return GetSecret("hmac")
      })

      if err != nil { // Token is expired
        c.Set("authorization", Anonymous)
        c.Set("user", "")
      }

      if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid { // Good token
        c.Set("authorization", Authenticated)
        c.Set("user", claims.Username)
      } else { // Tokens claims do not validate
        c.Set("authorization", Anonymous)
        c.Set("user", "")
      }
    }
  }
}

func Login(c *gin.Context, user string) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, GetAuthClaims(user))
  secret, err := GetSecret("hmac")
  if err != nil {
    panic(err)
  }
  signedTokenString, err := token.SignedString(secret)
  if err != nil {
    panic(err)
  }
  c.SetCookie("auth_jwt", signedTokenString, 3600, "/", "127.0.0.1", http.SameSiteNoneMode, false, false)
}

func IsValidAuth(c *gin.Context, store redis.Conn, user, pass string) bool {
  db_pass, err := redis.String(store.Do("HGET", fmt.Sprintf("user:%s", user), "password"))
  if err != nil {
    return false
  } else {
    return db_pass == pass
  }
}
