package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type call func(ctx context.Context) (*APIResponse, error)

func constructCall(req *APIRequest) call {
	return func(ctx context.Context) (*APIResponse, error) {
		request, err := http.NewRequestWithContext(ctx, req.Method, req.URI, req.BodyReader())
		if err != nil {
			return nil, fmt.Errorf("invalid request: %w", err)
		}
		request.Header = req.Headers

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("request failure: %w", err)
		}
		defer func() {
			_ = response.Body.Close()
		}()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		if len(responseData) > 0 {
			// TODO: Base this on Content Type header.
			responseData = formatJSON(responseData)
		}

		// TODO: This should not be added as a header. Instead, it should have
		// dedicated visualization.
		response.Header.Add("Status", response.Status)

		return &APIResponse{
			StatusCode: response.StatusCode,
			Body:       string(responseData),
			Headers:    response.Header,
		}, nil
	}
}

type APIRequest struct {
	Method  string
	URI     string
	Headers http.Header
	Body    string
}

func (r APIRequest) BodyReader() io.Reader {
	if r.Body == "" {
		return nil
	}
	return strings.NewReader(r.Body)
}

type APIResponse struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

func formatJSON(data []byte) []byte {
	var responseJSON any
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&responseJSON); err != nil {
		return data
	}

	var formattedData bytes.Buffer
	encoder := json.NewEncoder(&formattedData)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(responseJSON); err != nil {
		return data
	}

	return formattedData.Bytes()
}
