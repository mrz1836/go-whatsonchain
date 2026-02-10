package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// checkStatusCode returns an error if the last request returned an error HTTP status code.
// This prevents non-JSON error bodies (e.g. from 401, 429, 5xx) from being silently
// passed to json.Unmarshal, which would produce a confusing parse error.
//
// HTTP 404 is excluded because the API uses it for "resource not found" responses,
// which are handled by each endpoint's emptyErr parameter with domain-specific errors.
func checkStatusCode(c *Client, resp string) error {
	status := c.LastRequest().StatusCode
	if status == http.StatusOK || status == http.StatusNotFound {
		return nil
	}
	if len(resp) == 0 {
		return fmt.Errorf("%w: HTTP %d", ErrRequestFailed, status)
	}
	return fmt.Errorf("%w: HTTP %d: %s", ErrRequestFailed, status, resp)
}

// requestAndUnmarshal is a generic helper that performs a request and will unmarshal the response
// into a pointer to the specified type T
func requestAndUnmarshal[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) (*T, error) {
	resp, err := c.request(ctx, url, method, payload)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(c, resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, emptyErr
	}

	var result T
	if err = json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// requestAndUnmarshalSlice is a generic helper that performs a request and unmarshals the response
// into a slice of the specified type T
func requestAndUnmarshalSlice[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) ([]T, error) {
	resp, err := c.request(ctx, url, method, payload)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(c, resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, emptyErr
	}

	var result []T
	if err = json.Unmarshal([]byte(resp), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// requestString is a helper that performs a GET request and returns the raw string response
func requestString(ctx context.Context, c *Client, url string) (string, error) {
	return c.request(ctx, url, "GET", nil)
}
