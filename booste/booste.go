package booste

import (
	"encoding/json"
	"fmt"
)

// Run will call the inference pipeline on custom models with the use of a model key.
// It is a syncronous wrapper around the async Start and Check functions.
func Run(apiKey string, modelKey string, payloadIn interface{}, payloadOut interface{}) error {

	// Start the task
	taskID, err := Start(apiKey, modelKey, payloadIn)
	if err != nil {
		return err
	}

	// Poll check until done
	done := false
	for {
		done, err = Check(apiKey, taskID, payloadOut)
		if err != nil {
			return err
		}
		if done {
			break
		}
	}

	// payloadOut is now populated with returned data, so return no errors
	return nil
}

// The payload sent into the Start endpoint
type pStart struct {
	APIKey          string      `json:"apiKey"`
	ModelKey        string      `json:"modelKey"`
	ModelParameters interface{} `json:"modelParameters"` // send generic payloads.
}

// The response sent by the Start endpoint
type reStart struct {
	Status string `json:"status"`
	TaskID string `json:"taskID"`
}

// Start will start an async inference task and return a task ID
func Start(apiKey string, modelKey string, payloadIn interface{}) (string, error) {
	p := pStart{
		APIKey:          apiKey,
		ModelKey:        modelKey,
		ModelParameters: payloadIn, // name mismatch for backward compat to v1 backend, which expects modelParameters as json
	}

	re := reStart{}

	url := "http://localhost:5000/start"

	err := post(url, &p, &re)
	if err != nil {
		return "", err
	}

	if re.Status != "started" {
		return "", fmt.Errorf("inference task did not start")
	}

	if re.TaskID == "" {
		return "", fmt.Errorf("inference task returned an empty taskID")
	}

	return re.TaskID, nil
}

// The payload sent into the Check endpoint
type pCheck struct {
	APIKey string `json:"apiKey"`
	TaskID string `json:"taskID"`
}

// The response sent by the Check endpoint
type reCheck struct {
	Status string          `json:"status"`
	TaskID string          `json:"taskID"`
	Output json.RawMessage `json:"output"` // keep output raw for unmarshalling based off of payloadOut
}

// Check will check an existing async inference task and marshal the output into payloadOut if successful
// It returns (true, nil) if the task fully returned
// It returns (false, nil) if check call ran succesfully but the task has not yet returned into payloadOut
func Check(apiKey string, taskID string, payloadOut interface{}) (finished bool, err error) {

	p := pCheck{
		APIKey: apiKey,
		TaskID: taskID,
	}

	re := reCheck{}

	url := "http://localhost:5000/check"

	err = post(url, &p, &re)
	if err != nil {
		return false, err
	}

	// Handle running tasks, where the Check call ran without error, but the task is not done
	if re.Status == "started" || re.Status == "pending" || re.Status == "retrying" {
		return false, nil
	}

	// Handle failed tasks with an error
	if re.Status == "failed" {
		return false, fmt.Errorf("inference task failed")
	}

	// Else re.Status == "finished"
	err = json.Unmarshal(re.Output, payloadOut)
	if err != nil {
		return false, err
	}
	return true, nil
}
