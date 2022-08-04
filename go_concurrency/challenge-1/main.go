package main

import (
	"fmt"
	"sync"
)

var msg string

func updateMessage(s string, w *sync.WaitGroup) {
	defer w.Done()

	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

func main() {

	// challenge: modify this code so that the calls to updateMessage() on lines
	// 28, 30, and 33 run as goroutines, and implement wait groups so that
	// the program runs properly, and prints out three different messages.
	// Then, write a test for all three functions in this program: updateMessage(),
	// printMessage(), and main().

	// msg = "Hello, world!"

	// updateMessage("Hello, universe!")
	// printMessage()

	// updateMessage("Hello, cosmos!")
	// printMessage()

	// updateMessage("Hello, world!")
	// printMessage()

	var wg sync.WaitGroup

	xs := []string{"Hello, universe!", "Hello, cosmos!", "Hello, world!"}

	for _, s := range xs {
		wg.Add(1)
		go updateMessage(s, &wg)
		wg.Wait()
		printMessage()
	}
}
