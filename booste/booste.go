package booste

import (
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
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

		time.Sleep(time.Second)
	}

	// payloadOut is now populated with returned data, so return no errors
	return nil
}

// Start will start an async inference task and return a task ID.
func Start(apiKey string, modelKey string, payloadIn interface{}) (taskID string, err error) {

	data := inStartV1Data{
		APIKey:          apiKey,
		ModelKey:        modelKey,
		ModelParameters: payloadIn, // name mismatch for backward compat to v1 backend, which expects modelParameters as json
	}

	p := inStartV1{
		ID:      uuid.New().String(),
		Created: time.Now().Unix(),
		Data:    data,
	}

	re := outStartV1{}

	url := endpoint + "api/task/start/v1/"

	err = post(url, &p, &re)
	if err != nil {
		return "", err
	}

	if !re.Success {
		return "", fmt.Errorf(re.Message)
	}

	return re.Data.TaskID, nil
}

// Check will check the status of an existing async inference task. If the task has finished, the task's return values will be marshalled into payloadOut.
// The "done" boolean return value indicates if the requested async inference task has finished (true) or is still running (false).
func Check(apiKey string, taskID string, payloadOut interface{}) (done bool, err error) {

	data := inCheckV1Data{
		TaskID: taskID,
	}

	p := inCheckV1{
		ID:       uuid.New().String(),
		Created:  time.Now().Unix(),
		LongPoll: true,
		Data:     data,
	}

	re := outCheckV1{}

	url := endpoint + "api/task/check/v1/"

	err = post(url, &p, &re)
	if err != nil {
		return false, err
	}

	if !re.Success {
		return false, fmt.Errorf(re.Message)
	}

	if re.Data.TaskStatus == "Done" {
		err = json.Unmarshal(re.Data.TaskOut, payloadOut)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil

}
