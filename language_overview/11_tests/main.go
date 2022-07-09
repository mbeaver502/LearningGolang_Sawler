package main

import (
	"fmt"
	"log"
)

func main() {
	d, err := divide(100.0, 10.0)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("result:", d)
}

func divide(x, y float32) (float32, error) {
	if y == 0 {
		return 0, fmt.Errorf("divide by zero error: %v / %v", x, y)
	}

	return (x / y), nil
}
