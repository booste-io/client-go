package booste

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {

	// Define arbitrary struct to send in as payloadIn
	type pIn struct {
		A string
		B int
	}
	p := pIn{
		A: "This is an arbitrary struct",
		B: 1,
	}

	// Define the payloadOut to be returned
	// The mock server this is tested against returns json with a single key of "data"
	type reOut struct {
		Data string `json:"data"`
	}
	re := reOut{}

	err := Run("fakeAPIKey", "fakeModelKey", &p, &re)
	if err != nil {
		panic(err)
	}

	// re is now populated with results
	fmt.Printf("Out value: %+v\n", re)
}
