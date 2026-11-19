package whatsonchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"
)

// maxErrorBodyLength is the maximum number of response body bytes quoted in an
// error message. A response body can be up to maxResponseSize bytes (an HTML
// error page, for example), and error messages tend to end up in logs.
const maxErrorBodyLength = 512

// errorBodySnippet returns the response body, trimmed of surrounding whitespace,
// for quoting in an error message. A body longer than maxErrorBodyLength bytes is
// cut on a UTF-8 rune boundary and marked as truncated with its total size.
func errorBodySnippet(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) <= maxErrorBodyLength {
		return string(trimmed)
	}

	// Step back to the start of the rune that straddles the cut. A valid rune is
	// at most utf8.UTFMax bytes long, so for invalid UTF-8 (binary data) give up
	// after that many bytes rather than scanning the whole body.
	cut := maxErrorBodyLength
	for i := 0; i < utf8.UTFMax-1 && !utf8.RuneStart(trimmed[cut]); i++ {
		cut--
	}

	return fmt.Sprintf("%s... (truncated, %d bytes total)", trimmed[:cut], len(body))
}

// checkStatusCode returns an error if the HTTP status code indicates failure.
// This prevents non-JSON error bodies (e.g. from 401, 429, 5xx) from being silently
// passed to json.Unmarshal, which would produce a confusing parse error, and stops
// string endpoints from returning an error page as if it were data.
//
// HTTP 404 maps to notFoundErr when the endpoint has one:
//   - an empty body (how the API reports an unknown resource) returns notFoundErr itself
//   - any other body (e.g. "404 page not found" for an unknown route, usually a
//     misconfigured base URL) returns notFoundErr wrapped with the status and body
//
// errors.Is matches notFoundErr in both cases. Every other non-200 status, and a 404
// on an endpoint without a not-found sentinel (notFoundErr == nil), returns
// ErrRequestFailed with the status code and a truncated copy of the body.
func checkStatusCode(status int, body []byte, notFoundErr error) error {
	if status == http.StatusOK {
		return nil
	}

	snippet := errorBodySnippet(body)
	if status == http.StatusNotFound && notFoundErr != nil {
		if snippet == "" {
			return notFoundErr
		}
		return fmt.Errorf("%w: HTTP %d: %s", notFoundErr, status, snippet)
	}

	if snippet == "" {
		return fmt.Errorf("%w: HTTP %d", ErrRequestFailed, status)
	}
	return fmt.Errorf("%w: HTTP %d: %s", ErrRequestFailed, status, snippet)
}

// requestAndUnmarshal is a generic helper that performs a request and will unmarshal the response
// into a pointer to the specified type T.
//
// emptyErr is returned for an empty 200 response and for a 404 (see checkStatusCode).
// Pass nil for endpoints where a 404 is never a "not found" answer, such as bulk POST
// endpoints, so that a 404 is reported as ErrRequestFailed.
func requestAndUnmarshal[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) (*T, error) {
	resp, statusCode, err := c.request(ctx, url, method, payload)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(statusCode, resp, emptyErr); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, emptyErr
	}

	var result T
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// requestAndUnmarshalSlice is a generic helper that performs a request and unmarshals the response
// into a slice of the specified type T.
//
// emptyErr is handled the same way as in requestAndUnmarshal.
func requestAndUnmarshalSlice[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) ([]T, error) {
	resp, statusCode, err := c.request(ctx, url, method, payload)
	if err != nil {
		return nil, err
	}

	if err = checkStatusCode(statusCode, resp, emptyErr); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, emptyErr
	}

	var result []T
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// requestString is a helper that performs a GET request and returns the raw string response.
//
// notFoundErr is returned for a 404 (see checkStatusCode); pass nil to report a 404
// as ErrRequestFailed instead.
func requestString(ctx context.Context, c *Client, url string, notFoundErr error) (string, error) {
	resp, statusCode, err := c.request(ctx, url, http.MethodGet, nil)
	if err != nil {
		return "", err
	}

	if err = checkStatusCode(statusCode, resp, notFoundErr); err != nil {
		return "", err
	}

	return string(resp), nil
}
