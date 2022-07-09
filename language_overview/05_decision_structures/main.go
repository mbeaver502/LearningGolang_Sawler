package main

import "log"

func main() {
	isTrue := false

	if isTrue {
		log.Println("isTrue is", isTrue)
	} else {
		log.Println("isTrue is", isTrue)
	}

	cat := "dog"

	if cat == "cat" {
		log.Println("cat is", cat)
	} else {
		log.Println("cat is", cat)
	}

	num := 100
	isTrue = false

	if num > 99 && !isTrue {
		log.Println("hey presto")
	} else {
		log.Println("failure")
	}

	myVar := "lizard"

	switch myVar {
	case "dog":
		log.Println(myVar, "is dog")
	case "cat":
		log.Println(myVar, "is cat")
	default:
		log.Println(myVar, "is something else")
	}
}
