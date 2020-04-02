package main

import(
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

var(
  store redis.Conn
)


func main() {
  var err error
  store, err = redis.Dial("tcp", ":6379")
  if err != nil {
    panic(err)
  }
  defer store.Close()

  router := gin.Default()

  v1 := router.Group("/v1")
  {
    v1.POST("/users", statusCheck)                      // Create user account, create JWA, send to user
    v1.POST("/signin", statusCheck)                     // Create JWA, send to user
    v1.POST("/signout", statusCheck)                    // Destroy JWA

    v1.GET("/users/:uid", statusCheck)                  // Get user account
    v1.PUT("/users/:uid", statusCheck)                  // Update user account
    v1.DELETE("/users/:uid", statusCheck)               // Delete user account

    v1.GET("/users/:uid/accounts", statusCheck)         // Get accounts, without passwords
    v1.POST("/users/accounts", statusCheck)             // Create new account, generate and encrypt password

    v1.GET("/users/:uid/accounts/:aid", statusCheck)    // Get account, including password
    v1.PUT("/users/:uid/accounts/:aid", statusCheck)    // Update account
    v1.DELETE("/users/:uid/accounts/:aid", statusCheck) // Delete account
  }

  fmt.Println(redis.String(store.Do("Get", "myname")))

  router.Run(":6001")
}

func statusCheck(c *gin.Context) {
  c.String(http.StatusOK, "Healthy as a yak")
}
