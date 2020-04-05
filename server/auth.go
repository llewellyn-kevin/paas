package main

import(
  "crypto/sha256"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

const(
  // The name of the cookie on the client's device.
  CookieName = "auth_jwt"
  CookiePath = "/"
  CookieDomain = "127.0.0.1"
  // The name of the service, given in the JWT.
  ServiceName = "PaaS"
  // How long, in seconds, a token will stay valid.
  TokenLifetime = 60
  // If a request is made within this many seconds of the token expiring, a new
  // token will automatically be generated and sent.
  TokenRefreshWindow = 20

  // Codes used by the Authorization middleware to help controller method
  // handlers identify current authentication status of the client.
  Authenticated = "auth"
  Anonymous = "anoma"
)

// Reads the secret file for the given encoding method into a byte array that
// can be passed straight into the jwt method for signing tokens. This secret
// can be authomatically generated using a command from the makefile.
func GetSecret(method string) ([]byte, error) {
  return ioutil.ReadFile(fmt.Sprintf("jwt_secret.%s", method))
}

// AuthClaims is a model used to encode and parse the claims from the JWT, it 
// uses primarily the standard claims with one addition for the Username of the
// current user.
type AuthClaims struct {
  Username string `json:"usr"`
  jwt.StandardClaims
}

// GetAuthClaims takes the username of a user and returns an AuthClaims object
// that can be used to generate a JWT.
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

// Authorize is middleware used by Gin. Any route that uses the middleware 
// will look for a JWT in the request header, parse the claims, and check the
// signature. It will then set flags called "authorization" and "user" based
// on the validity of the given token. 
func Authorize() gin.HandlerFunc {
  return func(c *gin.Context) {
    now := time.Now().Unix()
    cookie, err := c.Cookie(CookieName)

    if err != nil { // No token has been set
      setAnonymous(c)
    } else {
      token, err := jwt.ParseWithClaims(cookie, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
        return GetSecret("hmac")
      })
      if err != nil { // Token is expired or could not be parsed
        setAnonymous(c)
      }

      if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid { // Good token
        log.Println(fmt.Sprintf("Token expires in %d s", claims.ExpiresAt - now))
        if claims.ExpiresAt - now <= TokenRefreshWindow { // Issue new token
          log.Println("need to make a new token")
          Logout(c)
          Login(c, claims.Username)
        }

        setAuthorized(c, claims.Username)
      } else { // Tokens claims could not be validated
        setAnonymous(c)
      }
    }
  }
}

// Hasher looks for any passwords in the query string and hashes them. 
func Hasher() gin.HandlerFunc {
  return func(c *gin.Context) {
    if c.Query("password") != "" {
      c.Set("password", sha256.Sum256([]byte(c.Query("password"))))
    }
  }
}

// setAnonymous sets flags to indicate to controller method handlers that
// the user who made the request has not supplied a valid authentication
// token. 
func setAnonymous(c *gin.Context) {
  c.Set("authorization", Anonymous)
  c.Set("user", "")
}

// setAuthorized sets flags to indicate to controller method handlers that
// the user who made the request is authorized, and what their username is.
func setAuthorized(c *gin.Context, username string) {
  c.Set("authorization", Authenticated)
  c.Set("user", username)
}

// Login takes a username to generate a JWT, and sets that token as a cookie on
// the client. This token will be sent by the client on all future requests and
// will be used to validate the requests by the Authorization middleware. This
// function merely creates a token for a given user. It does not validate that
// the user has supplied a valid password.
func Login(c *gin.Context, user string) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, GetAuthClaims(user))

  secret, err := GetSecret("hmac")
  if err != nil { // Secret file is not found
    panic(err)
  }

  signedTokenString, err := token.SignedString(secret)
  if err != nil {
    panic(err)
  }

  log.Println("Should create cookie now")
  c.SetCookie(CookieName, signedTokenString, TokenLifetime, CookiePath, CookieDomain, http.SameSiteNoneMode, false, false)
}

// Replaces the JWT on the client's device with an empty cookie.
func Logout(c *gin.Context) {
  c.SetCookie(CookieName, "", TokenLifetime, CookiePath, CookieDomain, http.SameSiteNoneMode, false, false)
}

// IsValidAuth takes a database connection, username, and password; and
// determines if the password is valid. Returns true if the user exists and
// the password matches. Returns false if the user could not be found or the
// password does not match.
func IsValidAuth(c *gin.Context, store redis.Conn, user, pass string) bool {
  db_pass, err := redis.String(store.Do("HGET", fmt.Sprintf("user:%s", user), "password"))
  if err != nil {
    return false
  } else {
    return db_pass == pass
  }
}
