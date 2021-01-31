package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	i, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		panic(err)
	}
	numberOfUuids := int(i)

	uuids := getUuids(numberOfUuids)

	for i, value := range uuids {
		fmt.Printf("%d - %s\n", i+1, value)
	}
}
