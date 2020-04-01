package main

type User struct {
	Id 			string
	FirstName 	string
	LastName 	string
	Email 		string
	Password 	string
	PubKey		string
	Accounts	[]string
}