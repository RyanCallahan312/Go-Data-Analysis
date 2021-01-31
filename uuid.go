package main

import (
	"github.com/google/uuid"
)

// getUuids generates a given amount of uuids and returns them as an array of strings
func getUuids(amount int) []string {

	if amount < 0 {
		return nil
	}

	uuids := make([]string, amount)

	for i := 0; i < amount; i++ {
		uuids[i] = uuid.New().String()
	}

	return uuids
}
