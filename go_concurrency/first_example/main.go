package main

import (
	"fmt"
	"log"
	"sync"
)

// main() is itself a goroutine
func main() {
	var wg sync.WaitGroup

	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"pi",
		"zeta",
		"eta",
		"theta",
		"epsilon",
	}

	wg.Add(len(words))
	for i, x := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, x), &wg)
	}
	wg.Wait()

	wg.Add(1)
	printSomething("This is the second thing to be printed", &wg)
}

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(s)
}
