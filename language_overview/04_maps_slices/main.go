package main

import (
	"log"
	"sort"
)

type User struct {
	FirstName string
	LastName  string
}

func main() {
	maps()
	slices()
}

func maps() {
	myMap := make(map[string]string)
	myMap2 := make(map[string]int)
	userMap := make(map[string]User)

	myMap["dog"] = "Fido"
	myMap["other_dog"] = "Bar"
	myMap["dog"] = "Foo"

	log.Println(myMap["dog"], myMap["other_dog"])

	myMap2["first"] = 1
	myMap2["second"] = 2

	log.Println(myMap2["first"], myMap2["second"])

	userMap["james"] = User{
		FirstName: "James",
		LastName:  "Bond",
	}

	log.Println(userMap["james"].FirstName, userMap["james"].LastName)
}

func slices() {
	var myStrings []string
	var myInts []int

	myStrings = append(myStrings, "hello", "world")
	myInts = append(myInts, 2, 1, 3)

	log.Println(myStrings)
	log.Println(myInts)

	sort.Ints(myInts)
	log.Println(myInts)

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] > numbers[j] })
	log.Println(numbers)
}
