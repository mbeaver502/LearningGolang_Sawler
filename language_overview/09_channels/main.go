package main

import (
	"log"

	"github.com/mbeaver502/LearningGolang_Sawler/language_overview/09_channels/helpers"
)

const NUMPOOL = 100

func CalculateValue(c chan int) {
	n := helpers.RandomNumber(NUMPOOL)
	c <- n // send a message through the channel
}

func main() {
	c := make(chan int)
	defer close(c) // close channel when function exits

	go CalculateValue(c) // launch a new goroutine

	log.Println(<-c) // receive a message through the channel
}
