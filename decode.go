package dvclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	headerEntityID = "OData-EntityId"
)

func extractIDFromCreateResponse(resp *http.Response) string {
	headerVal := resp.Header.Get(headerEntityID)
	vals := strings.Split(headerVal, "(")
	if len(vals) < 2 {
		return ""
	}

	id := strings.TrimSuffix(vals[1], ")")
	return id
}

func decodeRequestBody[T Validator](resp *http.Response) (T, error) {
	var row T
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&row); err != nil {
		return row, fmt.Errorf("decoding body: %w", err)
	}

	if err := row.Validate(); err != nil {
		return row, fmt.Errorf("validating body: %w", err)
	}
	return row, nil
}

func decodeRequestBodyList[T Validator](resp *http.Response) ([]T, error) {
	var listResponse tableListResponse[T]
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("decoding list body: %w", err)
	}

	rows := listResponse.Value
	if len(rows) == 0 {
		return rows, nil
	}

	if err := rows[0].Validate(); err != nil {
		return nil, fmt.Errorf("validating body: %w", err)
	}
	return rows, nil
}
