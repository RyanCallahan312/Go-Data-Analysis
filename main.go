package main

import (
	"fmt"
	"strings"

	"github.com/pborman/uuid"
)

func main() {
	var uuidWithHyphen uuid.UUID = uuid.NewRandom()
	var uuid string = strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	fmt.Println(uuid)
}
