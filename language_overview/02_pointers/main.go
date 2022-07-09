package main

import "log"

func main() {
	var myString = "Green"
	log.Println(myString)

	changeUsingParameter(&myString)
	log.Println(myString)
}

func changeUsingParameter(s *string) {
	*s = "Yellow"
}
