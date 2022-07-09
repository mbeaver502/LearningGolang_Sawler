package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	var whatToSay = "Goodbye, cruel world!"
	var i int = 16

	fmt.Println(whatToSay)
	fmt.Println("i is set to", i)

	whatWasSaid, whatWasAlsoSaid := saySomething()
	fmt.Println(whatWasSaid, whatWasAlsoSaid)
}

func saySomething() (string, string) {
	return "something", "something else"
}
