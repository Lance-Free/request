# Request

The `request` package provides a simple and flexible way to perform HTTP requests in Go. It abstracts away the
complexities of making HTTP requests and handling responses for RESTful APIs.

## Usage

### Making Requests

You can make various types of HTTP requests using the functions provided by this package:

- `Get(url string, options ...func(*RequestConfiguration))`: Perform an HTTP GET request.
- `Post(url string, options ...func(*RequestConfiguration))`: Perform an HTTP POST request.
- `Put(url string, options ...func(*RequestConfiguration))`: Perform an HTTP PUT request.
- `Delete(url string, options ...func(*RequestConfiguration))`: Perform an HTTP DELETE request.

### Request Configuration

You can customize your requests by providing options such as headers, parameters, and request body using functional
options:

- `WithHeader(key, value string)`: Add a header to the request.
- `WithHeaders(headers map[string]string)`: Add multiple headers to the request.
- `WithBody(body interface{})`: Set the request body.
- `WithParameter(key, value string)`: Add a query parameter to the request.
- `WithParameters(parameters map[string]string)`: Add multiple query parameters to the request.
- `WithCookie(key, value string)`: Add a cookie to the request.
- `WithCookies(cookies map[string]string)`: Add multiple cookies to the request.

### Handling Responses

The response body is automatically decoded into the desired type. If the request fails or the response code indicates an
error, a `request.Error` struct is returned with relevant information.

### Handling Errors

The `request` package introduces a custom error type to provide comprehensive error information during HTTP requests.
When an error occurs, a `request.Error` struct is returned, encapsulating the following details:

- `Code`: The HTTP status code from the server response (if available).
- `Body`: The raw response body from the server (if available).
- `Message`: A descriptive error message indicating the nature of the error.

To handle errors effectively, you can inspect the returned error value and examine the fields of the `request.Error`
struct. This allows you to implement specific error-handling logic based on the nature of the encountered error.

By leveraging the `Code` field of the `request.Error` struct, you can discern the type of error and tailor your
error-handling approach accordingly.

## Example

```go
package main

import (
	"fmt"
	"github.com/lance-free/request"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Make a GET request
	user, err := request.Get[User]("https://api.example.com/users/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("User:", user)
}
```

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.