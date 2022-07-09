package main

import (
	"log"

	"github.com/mbeaver502/LearningGolang_Sawler/language_overview/08_packages/helpers"
)

func main() {
	h := helpers.SomeType{
		TypeName:   "myType",
		TypeNumber: 42,
	}

	log.Println(h.TypeName, h.TypeNumber)
	helpers.DoSomething()
}
