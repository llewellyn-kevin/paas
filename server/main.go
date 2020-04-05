package main

import(
  "fmt"
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

var(
  store redis.Conn
)

func main() {
  var storeErr error
  store, storeErr = redis.Dial("tcp", ":6379")
  if storeErr != nil {
    panic(storeErr)
  }
  defer store.Close()

  router := gin.Default()
  router.Use(Authorize())

  v1 := router.Group("/v1")
  {
    v1.POST("/users", createUser)                       // Create user account, create JWA, send to user
    v1.POST("/signin", login)                           // Create JWA, send to user
    v1.POST("/signout", statusCheck)                    // Destroy JWA

    v1.GET("/users/:uid", statusCheck)                  // Get user account
    v1.PUT("/users/:uid", statusCheck)                  // Update user account
    v1.DELETE("/users/:uid", statusCheck)               // Delete user account

    v1.GET("/users/:uid/accounts", statusCheck)         // Get accounts, without passwords
    v1.POST("/users/:uid/accounts", statusCheck)        // Create new account, generate and encrypt password

    v1.GET("/users/:uid/accounts/:aid", statusCheck)    // Get account, including password
    v1.PUT("/users/:uid/accounts/:aid", statusCheck)    // Update account
    v1.DELETE("/users/:uid/accounts/:aid", statusCheck) // Delete account
  }

  router.Run(":6001")
}

func statusCheck(c *gin.Context) {
  log.Println(c.MustGet("authorization"))

  if(c.MustGet("authorization").(string) == Authenticated) {
    user := c.MustGet("user").(string)
    c.String(http.StatusOK, fmt.Sprintf("Hello %s, How are you? I am well.", user))
  } else {
    c.String(http.StatusOK, "Healthy as a yak")
  }
}

func createUser(c *gin.Context) {
  var newUser User
  if c.ShouldBindQuery(&newUser) != nil {
    log.Println(newUser.FirstName)
    log.Println(newUser.LastName)
    log.Println(newUser.Email)
    log.Println(newUser.Password)
  }

  exists, err := redis.Bool(store.Do("EXISTS", fmt.Sprintf("user:%s", newUser.Email)))
  if err != nil {
    panic(err)
  }

  if exists {
    c.String(http.StatusConflict, fmt.Sprintf("User with email %s already exists", newUser.Email))
  } else {
    // Create user in redis
    store.Do("HSET", fmt.Sprintf("user:%s", newUser.Email), "firstname", newUser.FirstName)
    store.Do("HSET", fmt.Sprintf("user:%s", newUser.Email), "lastname", newUser.LastName)
    store.Do("HSET", fmt.Sprintf("user:%s", newUser.Email), "password", newUser.Password)

    // Set cookie to track user
    Login(c, newUser.Email)

    // Send response
    c.String(http.StatusConflict, "User has been created")
  }
}

func login(c *gin.Context) {
  username := c.Query("username")
  password := c.Query("password")

  if IsValidAuth(c, store, username, password) {
    Login(c, username)
    c.String(http.StatusOK, "Logged in")
  } else {
    c.String(http.StatusUnauthorized, "Invalid username/password")
  }
}

