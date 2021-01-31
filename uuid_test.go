package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestGetUuids(t *testing.T) {

	amountsToTest := [8]int{0, 1, 10, -10, 100, -100, 1000, -1000}

	for _, amount := range amountsToTest {
		got := getUuids(amount)

		if amount < 0 {
			amount = 0
		}

		if len(got) != amount {
			t.Errorf("Got len(%d); Expected %d", len(got), amount)
		}

		for _, gotUuid := range got {
			_, err := uuid.Parse(gotUuid)

			if err != nil {

				t.Errorf("Got Invalid UUIDv4 UUID %s", gotUuid)
			}
		}

	}

}
