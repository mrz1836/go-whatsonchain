package whatsonchain

import (
	"context"
	"encoding/json"
)

// requestAndUnmarshal is a generic helper that performs a request and will unmarshal the response
// into a pointer to the specified type T
func requestAndUnmarshal[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) (*T, error) {
	resp, err := c.request(ctx, url, method, payload)
	if err != nil {
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
