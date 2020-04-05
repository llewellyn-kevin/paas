package main

import(
  "fmt"

  "github.com/gomodule/redigo/redis"
)

type User struct {
  FirstName	string    `form:"firstname"`
  LastName	string    `form:"lastname"`
  Email		  string    `form:"email"`
  Password	string    `form:"pass"`
}

// IsValid returns a bool that indicates whether or not all the values that
// have been unmarshelled into the struct are valid.
func (u User) IsValid() bool {
  if u.FirstName == "" {
    return false
  } else if u.LastName == "" {
    return false
  } else if u.Email == "" {
    return false
  } else if u.Password == "" {
    return false
  } else {
    return true
  }
}

// Save takes a User struct and puts it into the database. This will create
// a new entry to update an existing entry.
func (u User) Save() {
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "firstname", u.FirstName)
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "lastname", u.LastName)
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "password", u.Password)
}

// DoesExist returns true if the user exists in the database. Else false.
func (u User) DoesExist() (exists bool) {
  exists, err := redis.Bool(store.Do("EXISTS", fmt.Sprintf("user:%s", u.Email)))
  if err != nil {
    panic(err)
  }
  return
}
