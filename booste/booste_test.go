package booste

import (
	"fmt"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	// Define arbitrary struct to send in as payloadIn
	type pIn struct {
		String string `json:"string"`
		Length int    `json:"length"`
	}
	p := pIn{
		String: "I do not need food.",
		Length: 10,
	}

	// Define the responseOut to be returned
	type reOut struct {
		String string `json:"output"`
	}
	re := new(reOut)

	apiKey := os.Getenv("API_KEY")
	modelKey := os.Getenv("MODEL_KEY")

	err := Run(apiKey, modelKey, &p, &re)
	if err != nil {
		panic(err)
	}

	// re is now populated with results
	fmt.Printf("Out value: %+v\n", re.String)
}

func TestGPT2(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	modelSize := "gpt2"
	str := "hi this is a test for"
	var temperature float32 = 0.8
	length := 10
	windowMax := 50

	outStr, err := GPT2(apiKey, modelSize, str, length, temperature, windowMax)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// re is now populated with results
	fmt.Printf("Out value: %+v\n", outStr)
}

func TestGPT2XL(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	modelSize := "gpt2-xl"
	str := "hi this is a test for"
	var temperature float32 = 0.8
	length := 10
	windowMax := 50

	outStr, err := GPT2(apiKey, modelSize, str, length, temperature, windowMax)
	if err != nil {
		panic(err)
	}

	// re is now populated with results
	fmt.Printf("Out value: %+v\n", outStr)
}
