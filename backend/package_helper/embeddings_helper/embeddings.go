package embeddings_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Embed(ctx context.Context, text string) ([]float32, error) {
	// Prepare the request URL with the text query parameter
	encodedText := url.QueryEscape(text)
	url := fmt.Sprintf("http://localhost:8000/api/embeddings?text=%s", encodedText)

	// Send the GET request to the endpoint
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into a float32 array
	var embeddings []float32
	err = json.Unmarshal(body, &embeddings)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
