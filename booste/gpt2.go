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

// GPT2 will call the inference pipeline on gpt2 models.
func GPT2(apiKey string, modelSize string, str string, length int, temperature float32, windowMax int) (string, error) {
	if modelSize != "gpt2" && modelSize != "gpt2-xl" {
		return "", fmt.Errorf("did not pass valid modelSize argument of 'gpt2' or 'gpt2-xl'")
	}

	taskID, err := gpt2Start(apiKey, modelSize, str, length, temperature, windowMax)
	if err != nil {
		return "", err
	}

	var re struct {
		Output []string `json:"output"`
	}

	// Poll check until done
	done := false
	outStr := ""
	for {
		done, outStr, err = gpt2Check(apiKey, taskID, &re)
		if err != nil {
			return "", err
		}
		if done {
			break
		}

		time.Sleep(time.Second)
	}

	return outStr, nil
}

// The response sent by the Start endpoint
type reStart struct {
	Status string `json:"status"`
	TaskID string `json:"taskID"`
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

// The "done" boolean return value indicates if the requested async inference task has finished (true) or is still running (false).
func gpt2Check(apiKey string, taskID string, payloadOut interface{}) (done bool, outStr string, err error) {

	type inGPT2Check struct {
		TaskID string `json:"TaskID"`
		APIKey string `json:"apiKey"`
	}

	p := inGPT2Check{
		TaskID: taskID,
		APIKey: apiKey,
	}

	type outGPT2Check struct {
		Status string   `json:"Status"`
		Output []string `json:"Output"`
	}

	re := outGPT2Check{}

	url := endpoint + "inference/pretrained/gpt2/async/check/v2"

	err = post(url, &p, &re)
	if err != nil {
		return false, "", err
	}

	if re.Status == "Finished" {
		outStr := strings.Join(re.Output[:], " ")
		if err != nil {
			return false, "", err
		}
		return true, outStr, nil
	}

	return false, "", nil
}
