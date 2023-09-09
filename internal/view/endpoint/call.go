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

		var responseJSON any
		if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
			return nil, fmt.Errorf("json error: %w", err)
		}

		var responseBody bytes.Buffer
		encoder := json.NewEncoder(&responseBody)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(responseJSON); err != nil {
			return nil, fmt.Errorf("json error: %w", err)
		}

		return &APIResponse{
			StatusCode: response.StatusCode,
			Body:       responseBody.String(),
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
