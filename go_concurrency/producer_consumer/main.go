package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var (
	pizzasMade   int
	pizzasFailed int
	total        int
)

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func main() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	//print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in background
	go pizzeria(pizzaJob)

	// create and run a consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	// print out the end message
	color.Cyan("----------------------------------")
	color.Cyan("Done for the day!")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)
}

func (p *Producer) close() error {
	ch := make(chan error)
	p.quit <- ch

	// will return nil if channel closes successfully
	return <-ch
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	i := 0

	// run forever or until we receive a quit notification
	for {
		// try to make a pizza
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			// we tried to make a pizza -- send something to data channel
			case pizzaMaker.data <- *currentPizza:

			// we're ready to quit
			case quitChan := <-pizzaMaker.quit:
				// close channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		delay := rand.Intn(5) + 1
		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}
