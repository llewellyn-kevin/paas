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

func (u User) Save() {
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "firstname", u.FirstName)
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "lastname", u.LastName)
  store.Do("HSET", fmt.Sprintf("user:%s", u.Email), "password", u.Password)
}

func (u User) DoesExist() (exists bool) {
  exists, err := redis.Bool(store.Do("EXISTS", fmt.Sprintf("user:%s", u.Email)))
  if err != nil {
    panic(err)
  }
  return
}
