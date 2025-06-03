package dataverse

import (
	"context"
	"fmt"
	"net/http"
)

// Validator validates a struct that represents a table.
type Validator interface {
	Validate() error
}

// Table represnts a dataverse table. It wraps a [Client] and adds
// CRUD methods that decode and validate the JSON.
type Table[T Validator] struct {
	client       *Client
	expands      []string
	resourcePath string
}

type tableListResponse[T Validator] struct {
	Value []T `json:"value"`
}

// NewTable returns a table that wraps a [Client] with type-specific methods.
func NewTable[T Validator](client *Client, resource string) *Table[T] {
	return &Table[T]{client: client}
}

// Client returns the wrapped [Client].
func (t *Table[T]) Client() *Client {
	return t.client
}

// New creates a new row, it then calls GetByID and returns the row with any applied options.
func (t *Table[T]) New(ctx context.Context, newParams any, opts QueryOptions) (T, error) {
	var row T

	resp, err := t.client.MakeRequest(ctx, http.MethodPost, t.resourcePath, opts, newParams)
	if err != nil {
		return row, err
	}
	id := extractIDFromCreateResponse(resp)

	return t.GetByID(ctx, id, opts)

}

// GetByID returns a table row by its ID.
func (t *Table[T]) GetByID(ctx context.Context, id string, opts QueryOptions) (T, error) {
	var row T
	path := fmt.Sprintf("%s(%s)", t.resourcePath, id)

	if len(t.expands) > 0 {
		opts.AddExpand(t.expands...)
	}

	resp, err := t.client.MakeRequest(ctx, http.MethodGet, path, opts, nil)
	if err != nil {
		return row, err
	}
	return DecodeResponse[T](resp)
}

// GetByAltKey returns a table row by an alternate unique key.
func (t *Table[T]) GetByAltKey(ctx context.Context, key, value string, opts QueryOptions) (T, error) {
	var row T
	path := fmt.Sprintf("%s(%s='%s')", t.resourcePath, key, value)

	if len(t.expands) > 0 {
		opts.AddExpand(t.expands...)
	}

	resp, err := t.client.MakeRequest(ctx, http.MethodGet, path, opts, nil)
	if err != nil {
		return row, err
	}
	return DecodeResponse[T](resp)

}

// Lists returns a slice of table rows.
func (t *Table[T]) List(ctx context.Context, opts QueryOptions) ([]T, error) {

	if len(t.expands) > 0 {
		opts.AddExpand(t.expands...)
	}

	resp, err := t.client.MakeRequest(ctx, http.MethodGet, t.resourcePath, opts, nil)
	if err != nil {
		return nil, err
	}

	return DecodeResponseList[T](resp)
}

func (t *Table[T]) Update(ctx context.Context, id string, updateParams any, opts QueryOptions) (T, error) {
	var row T
	path := fmt.Sprintf("%s(%s)", t.resourcePath, id)

	_, err := t.client.MakeRequest(ctx, http.MethodPatch, path, opts, updateParams)
	if err != nil {
		return row, err
	}

	return t.GetByID(ctx, id, opts)
}

// Delete deletes the given row.
func (t *Table[T]) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s(%s)", t.resourcePath, id)
	_, err := t.client.MakeRequest(ctx, http.MethodDelete, path, QueryOptions{}, nil)
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultExpand sets the expand field for all operations.
func (t *Table[T]) SetDefaultExpand(fields []string) {
	t.expands = fields
}
