package main

import (
	"log"
	"time"
)

//var s string = "seven"

// Don't have all these vars laying around... organize them into a struct!
//var firstName string
//var lastName string
//var phoneNumber string
//var age int
//var birthDate time.Time

// There's no private/public concept in Go.
// In order to make things "public"-visible,
// we export them by capitalizing the name.
// Names with first char lowercase are unexported.
type User struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Age         int
	BirthDate   time.Time
}

func main() {
	var s2 = "six"

	s := "eight" // this is a different variable s, it's not the same as the package-level s!!!

	log.Println("s is", s)
	log.Println("s2 is", s2)

	saySomething("xxx")

	user := User{
		FirstName:   "James",
		LastName:    "Bond",
		PhoneNumber: "(123) 867-5309",
		Age:         32,
		BirthDate:   time.Time{},
	}

	log.Println(user.FirstName, user.LastName, user.PhoneNumber, user.Age, user.BirthDate)
	user.Print()
}

func saySomething(s string) (string, string) {
	log.Println("s inside func is", s) // this is a different s, not the same as the package-level s!!!

	return s, "World"
}

// structs can have functions attached to them
// Any var of type User can call this function
func (u *User) Print() {
	log.Println(u.FirstName, u.LastName)
}
