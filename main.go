package main

import (
	"fmt"
	"strconv"
)

func main() {

	fmt.Print("Enter how many uuid's you would like to generate: ")
	var input string
	fmt.Scanln(&input)

	numberOfUuids, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		panic(err)
	}

	uuids := getUuids(int(numberOfUuids))

	for i, value := range uuids {
		fmt.Printf("%d - %s\n", i+1, value)
	}
}
