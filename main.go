package main

import (
	"fmt"
	"os"

	_ "github.com/go-gota/gota"
)

type entry struct {
	drug_name string
	ai        string
	route     string
	applNO    string
	company   string
	date      string
	link      string
}

func main() {

	fmt.Println("Place holder")
	file, err := os.Open("Medication Guides.csv")
	if err != nil {
		fmt.Println("error while reading medication guides.csv", err)
	}
	file

}
