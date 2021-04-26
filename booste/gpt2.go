package booste

// THIS FILE WILL BE DEPRECIATED SOON

// With other language clients, we've ran GPT2 through booste.gpt2,
// and we're migrating to booste.run("gpt2") syntax to keep the clientside API
// simple and unchanging as more models are added.

import (
	"fmt"
	"strings"
	"time"
)

type pGPT2Start struct {
	String      string  `json:"string"`
	Length      int     `json:"length"`
	Temperature float32 `json:"temperature"`
	APIKey      string  `json:"apiKey"`
	ModelSize   string  `json:"modelSize"`
	WindowMax   int     `json:"windowMax"`
}

// The response sent by the Start endpoint
type reGPT2Start struct {
	Status string `json:"Status"`
	TaskID string `json:"TaskID"`
}

// GPT2 will call the inference pipeline on gpt2 models.
// It is a syncronous wrapper around the async GPT2Start and GPT2Check functions.
func GPT2(apiKey string, modelSize string, str string, length int, temperature float32, windowMax int) (string, error) {
	if modelSize != "gpt2" && modelSize != "gpt2-xl" {
		return "", fmt.Errorf("did not pass valid modelSize argument of 'gpt2' or 'gpt2-xl'")
	}

	taskID, err := gpt2Start(apiKey, modelSize, str, length, temperature, windowMax)
	if err != nil {
		return "", err
	}

	var re []string

	// Poll check until done
	done := false
	for {
		done, err = Check(apiKey, taskID, &re)
		if err != nil {
			return "", err
		}
		if done {
			break
		}

		time.Sleep(time.Second)
	}

	outStr := strings.Join(re[:], " ")

	return outStr, nil
}

// Start will start an async inference task and return a task ID.
func gpt2Start(apiKey string, modelSize string, str string, length int, temperature float32, windowMax int) (taskID string, err error) {

	p := pGPT2Start{
		String:      str,
		Length:      length,
		Temperature: temperature,
		APIKey:      apiKey,
		ModelSize:   modelSize, // Will be either gpt2 or gpt2-xl
		WindowMax:   windowMax,
	}

	re := reStart{}

	url := endpoint + "inference/pretrained/gpt2/async/start"

	fmt.Println("Posting to GPT2 endpoint", url)
	err = post(url, &p, &re)
	if err != nil {
		return "", err
	}

	if re.Status != "Started" {
		return "", fmt.Errorf("inference task did not start")
	}

	if re.TaskID == "" {
		return "", fmt.Errorf("inference task returned an empty taskID")
	}

	return re.TaskID, nil
}
