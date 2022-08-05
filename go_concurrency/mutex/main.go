package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func main() {
	msg = "Hello, world!"

	// data race! both goroutines access and write to same location
	wg.Add(2)
	go updateMessage("Hello, universe!")
	go updateMessage("Hello, cosmos!")
	wg.Wait()

	fmt.Println(msg)
}

func updateMessage(s string) {
	defer wg.Done()

	msg = s
}
