package booste

import (
	"fmt"
	"os"
	"testing"
)

func TestClip(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	prompts := []string{"A basketball", "A football"}
	images := []string{"https://upload.wikimedia.org/wikipedia/commons/7/7a/Basketball.png"}
	sims, err := CLIP(apiKey, prompts, images)
	if err != nil {
		// t.Fatal(err)
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", sims)
}
