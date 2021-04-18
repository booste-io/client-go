package booste

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func post(url string, p interface{}, re interface{}) error {

	// Marshal the payload
	jsonBytes, _ := json.Marshal(p)

	// Post it
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	// Catch non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booste returned status code %v", resp.StatusCode)
	}

	// Read body into bytes
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse returned bytes into the struct
	err = json.Unmarshal(bodyBytes, re)
	if err != nil {
		return err
	}
	return nil
}
