package main

type User struct {
  FirstName	string    `form:"firstname"`
  LastName	string    `form:"lastname"`
  Email		  string    `form:"email"`
  Password	string    `form:"pass"`
}
