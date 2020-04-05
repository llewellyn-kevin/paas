package main

import(
  "fmt"
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
)

type UserController struct {}

// Create takes the query string to instantiate a new user, and adds that user
// to persistent data.
func (u UserController) Create(c *gin.Context) {
  var newUser User
  if c.ShouldBindQuery(&newUser) != nil {
    log.Println(newUser.FirstName)
    log.Println(newUser.LastName)
    log.Println(newUser.Email)
    log.Println(newUser.Password)
  }

  if newUser.DoesExist() {
    c.String(http.StatusConflict, fmt.Sprintf("User with email %s already exists", newUser.Email))
  } else {
    newUser.Save()
    Login(c, newUser.Email)
    c.String(http.StatusConflict, "User has been created")
  }
}

// Show returns information on the requested user.
func (u UserController) Show(c *gin.Context) {

}

// Update changes the information for the given user in persistent data.
func (u UserController) Update(c *gin.Context) {

}

// Destroy finds the given user in persistent data, and deletes it.
func (u UserController) Destroy(c *gin.Context) {

}

// Login takes a username and password from the query string and passes them
// onto the authentication service to generate a response
func (u UserController) Login(c *gin.Context) {
  username := c.Query("username")
  password := c.Query("password")

  if IsValidAuth(c, store, username, password) {
    Login(c, username)
    c.String(http.StatusOK, "Logged in")
  } else {
    c.String(http.StatusUnauthorized, "Invalid username/password")
  }
}

// Logout uses the auth service to destroy the current JWT
func (u UserController) Logout(c *gin.Context) {
  Logout(c)
}
