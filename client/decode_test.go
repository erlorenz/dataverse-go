package dv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type something struct {
	ID    string `json:"somethingid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s something) Validate() error {
	if s.ID == "" || s.Name == "" || s.Email == "" {
		return errors.New("validate error")
	}
	return nil
}

func TestExtractID(t *testing.T) {
	randomID := "RANDOM_ID"
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Add(headerEntityID, fmt.Sprintf("https://company.api.crm.dynamics.com/api/data/v9.2/entities(%s)", randomID))

	want := randomID
	got := extractIDFromCreateResponse(resp)

	if want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestDecodeBody(t *testing.T) {
	want := something{"SOME_ID", "John Doe", "jdoe@example.com"}
	b, _ := json.Marshal(want)

	resp := &http.Response{}
	resp.Body = io.NopCloser(bytes.NewReader(b))

	got, err := decodeRequestBody[something](resp)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestDecodeListBody(t *testing.T) {

	want := []something{
		{"SOME_ID", "John Doe", "jdoe@example.com"},
		{"SOME_OTHER_ID", "Brett Fox", "bfox@example.com"},
	}
	responseBody := tableListResponse[something]{Value: want}
	respBodyBytes, _ := json.Marshal(responseBody)

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(respBodyBytes))}

	got, err := decodeRequestBodyList[something](resp)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
