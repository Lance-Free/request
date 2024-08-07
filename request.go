package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Method represents an HTTP method.
type Method string

const (
	getMethod    Method = "GET"
	postMethod   Method = "POST"
	putMethod    Method = "PUT"
	deleteMethod Method = "DELETE"
)

// RequestConfiguration represents the configuration for an HTTP request.
type RequestConfiguration struct {
	Headers, Parameters, Cookies map[string]string
	Body                         []byte
}

// Error represents an error returned by the application.
type Error struct {
	Code    int
	Body    []byte
	Message string
}

// Error returns the error message associated with the Error struct.
func (e *Error) Error() string {
	return e.Message
}

// do is a generic function that performs an HTTP request with the specified method, URL, and request options.
// It returns the response body decoded into the type T and an error if the request fails or the response code indicates an error.
func do[T any](method Method, url string, options ...func(*RequestConfiguration)) (T, *Error) {
	requestConfig := RequestConfiguration{
		Headers:    make(map[string]string),
		Parameters: make(map[string]string),
		Cookies:    make(map[string]string),
	}
	for _, option := range options {
		option(&requestConfig)
	}

	var result T
	client := http.DefaultClient
	request, err := http.NewRequest(string(method), url, io.NopCloser(bytes.NewReader(requestConfig.Body)))
	if err != nil {
		return result, &Error{Message: "failed to create request"}
	}

	for key, value := range requestConfig.Headers {
		request.Header.Set(key, value)
	}

	values := request.URL.Query()
	for key, value := range requestConfig.Parameters {
		values.Add(key, value)
	}
	request.URL.RawQuery = values.Encode()

	for key, value := range requestConfig.Cookies {
		request.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	resp, err := client.Do(request)
	if err != nil {
		return result, &Error{Message: "failed to send request"}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if !isSuccessfulCode(resp.StatusCode) {
		responseBytes, _ := io.ReadAll(resp.Body)
		return result, &Error{
			Code:    resp.StatusCode,
			Body:    responseBytes,
			Message: fmt.Sprintf("status code does not indicate success: %d", resp.StatusCode),
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			return result, &Error{Message: fmt.Sprintf("failed to decode field \"%s\"", unmarshalErr.Field)}
		}
		return result, &Error{Message: "failed to decode response"}
	}

	return result, nil
}

// Get is a generic function that performs an HTTP GET request with the specified URL and request options.
// It returns the response body decoded into the type T and an error if the request fails or the response code indicates an error.
func Get[T any](url string, options ...func(*RequestConfiguration)) (T, *Error) {
	return do[T](getMethod, url, options...)
}

// Post is a generic function that performs an HTTP GET request with the specified URL and request options.
// It returns the response body decoded into the type T and an error if the request fails or the response code indicates an error.
func Post[T any](url string, options ...func(*RequestConfiguration)) (T, *Error) {
	return do[T](postMethod, url, options...)
}

// Put is a generic function that performs an HTTP GET request with the specified URL and request options.
// It returns the response body decoded into the type T and an error if the request fails or the response code indicates an error.
func Put[T any](url string, options ...func(*RequestConfiguration)) (T, *Error) {
	return do[T](putMethod, url, options...)
}

// Delete is a generic function that performs an HTTP GET request with the specified URL and request options.
// It returns the response body decoded into the type T and an error if the request fails or the response code indicates an error.
func Delete[T any](url string, options ...func(*RequestConfiguration)) (T, *Error) {
	return do[T](deleteMethod, url, options...)
}

// isSuccessfulCode checks if the given code falls within the range of successful HTTP status codes.
// It returns true if the code is greater than or equal to 200 and less than 300; otherwise, it returns false.
func isSuccessfulCode(code int) bool {
	return code >= 200 && code < 300
}

// WithHeader adds a header to the request.
// It takes a key string and a value string as parameters.
// It returns a function that takes a pointer to a RequestConfiguration struct as a parameter.
// This function updates the Headers field of the RequestConfiguration struct by adding the key-value pair.
func WithHeader(key, value string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		r.Headers[key] = value
	}
}

// WithHeaders is a function that takes a map of headers and returns another function that
// modifies the headers of a RequestConfiguration object by adding the key-value pairs from the input map.
func WithHeaders(headers map[string]string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		for key, value := range headers {
			r.Headers[key] = value
		}
	}
}

// WithAccept is a function that returns an option function for setting the "Accept" header in an HTTP request.
// The option function adds the "Accept" header with the value "application/json" to the RequestConfiguration structure.
func WithAccept() func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		r.Headers["Accept"] = "application/json"
	}
}

// WithBody is a higher-order function that takes a value of any type as its argument and returns
// a function that takes a pointer to a RequestConfiguration object as its argument. The returned function sets
// the Body field of the provided RequestConfiguration object with the JSON representation of the input value.
// If the marshaling of the input value fails, the Body field will remain unchanged.
func WithBody[T any](body T) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return
		}

		r.Body = jsonBody
		r.Headers["Content-Type"] = "application/json"
	}
}

// WithParameter is a higher-order function that returns a function which
// sets a parameter with a given key and value on a RequestConfiguration object.
// The returned function modifies the Parameters map of the provided RequestConfiguration by
// adding or updating the entry with the given key and value.
func WithParameter(key, value string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		r.Parameters[key] = value
	}
}

// WithParameters accepts a map of parameters and returns a closure that
// modifies a given *RequestConfiguration by adding all the key-value pairs from the
// parameters map to the RequestConfiguration's Parameters field.
func WithParameters(parameters map[string]string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		for key, value := range parameters {
			r.Parameters[key] = value
		}
	}
}

// WithCookie is a higher-order function that returns a function which
// sets a cookie with a given key and value on a RequestConfiguration object.
// The returned function modifies the Cookies map of the provided RequestConfiguration by
// adding or updating the entry with the given key and value.
func WithCookie(key, value string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		r.Cookies[key] = value
	}
}

// WithCookies accepts a map of cookies and returns a closure that
// modifies a given *RequestConfiguration by adding all the key-value pairs from the
// cookies map to the RequestConfiguration's Cookies field.
func WithCookies(cookies map[string]string) func(*RequestConfiguration) {
	return func(r *RequestConfiguration) {
		for key, value := range cookies {
			r.Cookies[key] = value
		}
	}
}
