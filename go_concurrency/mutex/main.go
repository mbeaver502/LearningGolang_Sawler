package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func main() {
	msg = "Hello, world!"

	var mutex sync.Mutex

	// data race! both goroutines access and write to same location
	wg.Add(2)
	go updateMessage("Hello, universe!", &mutex)
	go updateMessage("Hello, cosmos!", &mutex)
	wg.Wait()

	fmt.Println(msg)
}

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()

	// mutex allows us to safely access and write msg
	m.Lock()
	msg = s
	m.Unlock()
}
