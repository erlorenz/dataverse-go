package dvclient

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

const (
	OrderByAsc  = "asc"
	OrderByDesc = "desc"
)

// QueryOptions are optional fields that apply to all table operations.
type QueryOptions struct {
	Select  []string
	Expand  []string
	OrderBy string
	Filter  []string
	Top     int
}

// SetOrderBy is a convenience method to set the OrderBy field.
func (qo *QueryOptions) SetOrderBy(field, direction string) {
	qo.OrderBy = fmt.Sprintf("%s %s", field, direction)
}

// AddExpand is a convenience method to set multiple fields to expand.
// This is in addition to the default expands.
func (qo *QueryOptions) AddExpand(fields ...string) {
	qo.Expand = slices.Concat(qo.Expand, fields)
}

// AddFilter adds a filter joined with " and ".
func (qo *QueryOptions) AddFilter(filters ...string) {
	qo.Filter = slices.Concat(qo.Filter, filters)
}

// AddSelect is a convenience method add fields to the select param.
func (qo *QueryOptions) AddSelect(fields ...string) {
	qo.Select = slices.Concat(qo.Select, fields)
}

func (qo *QueryOptions) ToParams() url.Values {
	params := url.Values{}

	switch {
	case len(qo.Expand) > 0:
		params.Set("$expand", strings.Join(qo.Expand, ","))
		fallthrough
	case qo.OrderBy != "":
		params.Set("$orderby", qo.OrderBy)
		fallthrough
	case len(qo.Select) > 0:
		params.Set("$select", strings.Join(qo.Select, ","))
		fallthrough
	case len(qo.Filter) > 0:
		params.Set("$filter", strings.Join(qo.Filter, " and "))
	}

	return params
}
