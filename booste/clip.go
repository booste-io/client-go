package booste

// THIS FILE WILL BE DEPRECIATED SOON

// With other language clients, we've ran CLIP through booste.clip,
// and we're migrating to booste.run("clip") syntax to keep the clientside API
// simple and unchanging as more models are added.

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
)

type pCLIP struct {
	Prompt string `json:"prompt"`
	Image  string `json:"image"`
	APIKey string `json:"apiKey"`
	IsURL  bool   `json:"isUrl"`
}

// The response sent by the Start endpoint
type reCLIP struct {
	Similarity string `json:"similarity"`
	Message    string `json:"message"`
}

// CLIP will call the inference pipeline on a CLIP zero-shot image classifier model and return a nested map of similarities.
func CLIP(apiKey string, prompts []string, images []string) (map[string]map[string]map[string]float32, error) {

	// check lengths of inputs
	if len(prompts) < 1 {
		return nil, fmt.Errorf("no prompts were provided")
	}
	if len(images) < 1 {
		return nil, fmt.Errorf("no images were provided")
	}

	// validate urls for images
	for _, image := range images {
		_, err := url.ParseRequestURI(image)
		if err != nil {
			return nil, fmt.Errorf("image %v is not a valid url", image)
		}
	}

	// build return map to be filled out by parallel goroutines
	var sims = map[string]map[string]map[string]float32{}
	for _, prompt := range prompts {
		sims[prompt] = map[string]map[string]float32{}
		for _, image := range images {
			sims[prompt][image] = map[string]float32{}
			sims[prompt][image]["similarity"] = 0.0 // zero init
		}
	}

	sem := make(chan (map[string]map[string]map[string]float32), 1)
	sem <- sims

	errChan := make(chan error, len(prompts)*len(images))

	var wg sync.WaitGroup
	for _, prompt := range prompts {
		for _, image := range images {
			wg.Add(1)
			go callCLIP(apiKey, prompt, image, sem, errChan, &wg)
		}
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	sims = <-sem

	return sims, nil
}

// Start will start an async inference task and return a task ID.
func callCLIP(
	apiKey string,
	prompt string,
	image string,
	sem chan map[string]map[string]map[string]float32,
	errChan chan error,
	wg *sync.WaitGroup) {

	defer wg.Done()

	p := pCLIP{
		Prompt: prompt,
		Image:  image,
		APIKey: apiKey,
		IsURL:  true,
	}

	re := reCLIP{}

	url := "https://7rq1vzhvxj.execute-api.us-west-1.amazonaws.com/Prod/infer/"
	err := post(url, &p, &re)
	if err != nil {
		errChan <- err
		return
	}

	sim, err := strconv.ParseFloat(re.Similarity, 32)
	if err != nil {
		errChan <- err
		return
	}

	// pull similarity map from semaphore as mutex
	sims := <-sem
	sims[prompt][image]["similarity"] = float32(sim)
	sem <- sims
	errChan <- nil
}
