package main

import "log"

type Animal interface {
	Says() string
	NumberOfLegs() int
}

type Dog struct {
	Name  string
	Breed string
}

// best practice is to prefer using receivers of form (t *T)
func (d *Dog) Says() string {
	return "woof"
}

func (d *Dog) NumberOfLegs() int {
	return 4
}

type Gorilla struct {
	Name          string
	Color         string
	NumberOfTeeth int
}

func (g *Gorilla) Says() string {
	return "ooga booga"
}

func (g *Gorilla) NumberOfLegs() int {
	return 2
}

func main() {
	dog := Dog{
		Name:  "Samson",
		Breed: "German Shepherd",
	}

	gorilla := Gorilla{
		Name:          "Suzy",
		Color:         "Gray",
		NumberOfTeeth: 20,
	}

	PrintInfo(&dog)
	PrintInfo(&gorilla)
}

func PrintInfo(a Animal) {
	log.Println(a.Says(), a.NumberOfLegs())
}
