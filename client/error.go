package dataverse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	ErrorTypeConnection = "connection"
	ErrorTypeAPI        = "api"
	ErrorTypeAuth       = "auth"
)

var ErrInvalidAPIError = errors.New("invalid api error format")

type apiErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type APIError struct {
	Path       string
	StatusCode int
	Code       string
	Message    string
}

func (err APIError) Error() string {
	return fmt.Sprintf("%s  %s", err.Code, err.Message)
}

func handleAPIError(resp *http.Response) error {

	contentType := resp.Header.Get("Content-Type")

	if contentType != contentTypeJSON {
		return fmt.Errorf("%w: status %d, content-type '%s'", ErrInvalidAPIError, resp.StatusCode, contentType)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("%w: cannot read body", ErrInvalidAPIError)
	}

	var errResponse apiErrorResponse
	if err := json.Unmarshal(bodyBytes, &errResponse); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAPIError, err)
	}

	if errResponse.Error.Code == "" || errResponse.Error.Message == "" {
		return fmt.Errorf("%w: code or message is empty", ErrInvalidAPIError)
	}

	return APIError{
		Path:       resp.Request.URL.Path,
		StatusCode: resp.StatusCode,
		Code:       errResponse.Error.Code,
		Message:    errResponse.Error.Message,
	}

}
