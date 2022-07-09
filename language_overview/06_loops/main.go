package main

import "log"

type User struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

func main() {
	for i := 0; i < 10; i++ {
		log.Println(i)
	}

	animals := []string{"dog", "fish", "horse", "cow", "cat"}

	for _, animal := range animals {
		log.Println(animal)
	}

	m := make(map[string]string)
	m["dog"] = "Fido"
	m["cat"] = "Fluffy"

	for k, v := range m {
		log.Println(k, v)
	}

	s := "hello world!"

	// a string is a slice of bytes
	// string is an immutable type!
	for _, c := range s {
		log.Println(string(c))
	}

	users := []User{
		{
			FirstName: "James",
			LastName:  "Bond",
			Email:     "jbond@example.com",
			Age:       32,
		},
		{
			FirstName: "David",
			LastName:  "Smith",
			Email:     "davidsmith@example.com",
			Age:       35,
		},
		{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alicej@example.com",
			Age:       27,
		},
	}

	for _, user := range users {
		log.Println(user.FirstName, user.LastName, user.Email, user.Age)
	}
}
