package main

import(
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
)

type UserController struct {}

// Create takes the query string to instantiate a new user, and adds that user
// to persistent data.
func (u UserController) Create(c *gin.Context) {
  var newUser User
  newUser.FirstName = c.Query("firstname")
  newUser.LastName = c.Query("lastname")
  newUser.Email = c.Query("email")
  pass, _ := c.Get("password")
  newUser.Password = fmt.Sprintf("%v", pass)

  if !newUser.IsValid() {
    c.String(http.StatusBadRequest, "Invalid input")
  } else if newUser.DoesExist() {
    c.String(http.StatusConflict, fmt.Sprintf("User with email %s already exists", newUser.Email))
  } else {
    newUser.Save()
    Login(c, newUser.Email)
    c.String(http.StatusConflict, "User has been created")
  }
}

// Show returns information on the requested user.
func (u UserController) Show(c *gin.Context) {
  auth, _ := c.Get("authorization")
  if(auth != Authenticated) {
    c.String(http.StatusUnauthorized, "Not authorized for that action.")
  } else {
    user, _ := c.Get("user")
    c.String(http.StatusOK, fmt.Sprintf("Current user is: %s", user))
  }
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
  password, exists := c.Get("password")

  if IsValidAuth(c, store, username, fmt.Sprintf("%v", password)) && exists {
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
