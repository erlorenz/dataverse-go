package dv

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestError_APIError(t *testing.T) {

	fakeErrorBytes := []byte(`{"error": {"code": "not found","message":"not found message"}}`)

	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: 404,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader(fakeErrorBytes))}

	fullPath, _ := url.Parse("https://example.com/resourcepath")
	resp.Request.URL = fullPath
	resp.Header.Set("Content-Type", contentTypeJSON)

	err := handleAPIError(resp)
	var apiError APIError
	if !errors.As(err, &apiError) {
		t.Errorf("wanted %T, got %s", APIError{}, err)
	}

	want := APIError{Code: "not found", Message: "not found message", Path: "/resourcepath", StatusCode: 404}
	if !cmp.Equal(want, apiError) {
		t.Error(cmp.Diff(want, apiError))
	}
}

func TestError_InvalidContentType(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Content-Type", "text/html")

	err := handleAPIError(resp)
	if !errors.Is(err, ErrInvalidAPIError) {
		t.Errorf("expected ErrInvalidError, got %s", err)
	}

}

func TestError_InvalidErrorMissingFields(t *testing.T) {

	var invalidErrorBytes = []byte(`{"error" {"wrong": "field"}}`)
	resp := &http.Response{
		Request:    &http.Request{},
		Header:     http.Header{},
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader(invalidErrorBytes)),
	}
	resp.Header.Set("Content-Type", contentTypeJSON)

	err := handleAPIError(resp)
	if !errors.Is(err, ErrInvalidAPIError) {
		t.Errorf("expected ErrInvalidrror, got %s", err)
	}

}

func TestError_InvalidErrorUnmarshalError(t *testing.T) {

	var invalidErrorBytes = []byte(`{"error": 2`)
	resp := &http.Response{
		Request:    &http.Request{},
		Header:     http.Header{},
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader(invalidErrorBytes)),
	}
	resp.Header.Set("Content-Type", contentTypeJSON)

	err := handleAPIError(resp)
	if !errors.Is(err, ErrInvalidAPIError) {
		t.Errorf("expected ErrInvalidrror, got %s", err)
	}

}
