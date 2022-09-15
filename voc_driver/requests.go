package vocdriver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Evaluation Tools
var httpSuccessCodes = [...]int{
	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNoContent,
}

// RequestService is a handle to HTTP request operations
type RequestService struct {
	client *Client
}

// SuccessfulHTTPRequest returns true if the given Response's StatusCode
// is one of `[...]int{200, 201, 202, 204}`; otherwise returns false
// Taiga does not return status codes other than above stated
func SuccessfulHTTPRequest(response *http.Response) bool {
	for _, code := range httpSuccessCodes {
		if response.StatusCode == code {
			return true
		}
	}
	return false
}

// Get a handler for composing a new HTTP GET request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Get(url string, responseBody interface{}) (*http.Response, error) {
	return newRawRequest("GET", s.client, responseBody, url, nil)
}

// Head a handler for composing a new HTTP HEAD request
func (s *RequestService) Head() {
	panic("HEAD requests are not implemented")
}

// Post a handler for composing a new HTTP POST request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Post(url string, payload interface{}, responseBody interface{}) (*http.Response, error) {
	return newRawRequest("POST", s.client, responseBody, url, payload)
}

// Put a handler for composing a new HTTP PUT request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Put(url string, payload interface{}, responseBody interface{}) (*http.Response, error) {
	return newRawRequest("PUT", s.client, responseBody, url, payload)
}

// Patch a handler for composing a new HTTP PATCH request
//
//   - URL must be an absolute (full) URL to the desired endpoint
//   - Payload must be a pointer to a complete struct which will be sent to Taiga
//   - ResponseBody must be a pointer to a struct representing the fields returned by Taiga
func (s *RequestService) Patch(url string, payload interface{}, responseBody interface{}) (*http.Response, error) {
	return newRawRequest("PATCH", s.client, responseBody, url, payload)
}

// Delete a handler for composing a new HTTP DELETE request
//
//   - URL must be an absolute (full) URL to the desired endpoint
func (s *RequestService) Delete(url string) (*http.Response, error) {
	return newRawRequest("DELETE", s.client, nil, url, nil)
}

// Connect a handler for composing a new HTTP CONNECT request
func (s *RequestService) Connect() {
	panic("CONNECT requests are not implemented")
}

// Options a handler for composing a new HTTP OPTIONS request
func (s *RequestService) Options() {
	panic("OPTIONS requests are not implemented")
}

// Trace a handler for composing a new HTTP TRACE request
func (s *RequestService) Trace() {
	panic("TRACE requests are not implemented")
}

func newRawRequest(requestType string, c *Client, responseBody interface{}, url string, payload interface{}) (*http.Response, error) {
	// New RAW request
	var request *http.Request
	var err error

	switch {
	case payload == nil:
		request, err = http.NewRequest(requestType, url, nil)
		if err != nil {
			return nil, err
		}

	case payload != nil:
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(requestType, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("failed to build request because the received payload could not be processed")
	}

	// Load Headers
	c.loadHeaders(request)

	// Execute request
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Evaluate response status code
	if SuccessfulHTTPRequest(resp) {
		if resp.StatusCode == http.StatusNoContent { //  There's no body returned for 204 responses
			return resp, nil
		}
		// We expect content so convert response JSON string to struct
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&responseBody)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	rawResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf(string(rawResponseBody))
}
