package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts an API key from headers of an HTTP request
// Ex:
// Authorization: ApiKey { some_api_key }
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("No authentication info found")
	}
	vals := strings.Split(val, " ") // split result on spaces
	if len(vals) != 2 {
		return "", errors.New("Malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("First part of auth header malformed")
	}
	return vals[1], nil
}
