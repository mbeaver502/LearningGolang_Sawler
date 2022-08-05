package main

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

	//print out a message

	// create a producer

	// run the producer in background

	// create and run a consumer

	// print out the end message
}
