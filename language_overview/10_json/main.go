package main

import (
	"encoding/json"
	"log"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HairColor string `json:"hair_color"`
	HasDog    bool   `json:"has_dog"`
}

func main() {
	myJson := `
	[
		{
			"first_name": "Clark",
			"last_name": "Kent",
			"hair_color": "black",
			"has_dog": true
		},
		{
			"first_name": "Bruce",
			"last_name": "Wayne",
			"hair_color": "black",
			"has_dog": false
		}
	]`

	people := []Person{}
	err := json.Unmarshal([]byte(myJson), &people)

	if err != nil {
		log.Println("error:", err)
	}

	for _, person := range people {
		log.Println(person)
	}

	people = append(people,
		Person{
			FirstName: "Wally",
			LastName:  "West",
			HairColor: "red",
			HasDog:    true,
		},
		Person{
			FirstName: "Diana",
			LastName:  "Prince",
			HairColor: "black",
			HasDog:    false,
		})

	// json.MarshalIndent() is handy for pretty-printing, but prefer json.Marshal() in real world
	newJson, err := json.MarshalIndent(people, "", "   ")

	if err != nil {
		log.Println("error:", err)
	}

	log.Println(string(newJson))
}
