package testhelpers

import (
	"fmt"
	"net/http"
)

// Intp returns a pointer to the given integer value.
func Intp(i int) *int {
	return &i
}

// MustNewRequest creates a new http get request or panics if it can't
func MustNewRequest(rawurl string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
			panic(fmt.Sprintf("failed to create HTTP Get Request for '%s': %s", rawurl, err))
		}
	return request
}