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
  var storeErr error
  store, storeErr = redis.Dial("tcp", ":6379")
  if storeErr != nil {
    panic(storeErr)
  }
  defer store.Close()

  router := gin.Default()
  router.Use(Authorize())
  router.Use(Hasher())

  userController := UserController{}
  accountController := AccountController{}

  v1 := router.Group("/v1")
  {
    v1.POST("/users", userController.Create)                          // Create user account, create JWA, send to user
    v1.POST("/signin", userController.Login)                          // Create JWA, send to user
    v1.POST("/signout", userController.Logout)                        // Destroy JWA

    v1.GET("/users/:uid", userController.Show)                        // Get user account
    v1.PUT("/users/:uid", userController.Update)                      // Update user account
    v1.DELETE("/users/:uid", userController.Destroy)                  // Delete user account

    v1.GET("/users/:uid/accounts", accountController.Index)           // Get accounts, without passwords
    v1.POST("/users/:uid/accounts", accountController.Create)         // Create new account, generate and encrypt password

    v1.GET("/users/:uid/accounts/:aid", accountController.Show)       // Get account, including password
    v1.PUT("/users/:uid/accounts/:aid", accountController.Update)     // Update account
    v1.DELETE("/users/:uid/accounts/:aid", accountController.Destroy) // Delete account
  }

  router.Run(":6001")
}

func statusCheck(c *gin.Context) {
  if(c.MustGet("authorization").(string) == Authenticated) {
    user := c.MustGet("user").(string)
    c.String(http.StatusOK, fmt.Sprintf("Hello %s, How are you? I am well.", user))
  } else {
    c.String(http.StatusOK, "Healthy as a yak")
  }
}

