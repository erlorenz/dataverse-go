package dataverse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOptions_AddSelect(t *testing.T) {
	opts := QueryOptions{}
	opts.AddSelect("field1", "field2")
	opts.AddSelect("field3")

	want := []string{"field1", "field2", "field3"}
	got := opts.Select
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	params := opts.ToParams()
	want2 := "field1,field2,field3"
	got2 := params.Get("$select")

	if !cmp.Equal(want2, got2) {
		t.Error(cmp.Diff(want2, got2))
	}

}

func TestOptions_AddExpand(t *testing.T) {
	opts := QueryOptions{}
	opts.AddExpand("field1", "field2")
	opts.AddExpand("field3")

	want := []string{"field1", "field2", "field3"}
	got := opts.Expand
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	params := opts.ToParams()
	want2 := "field1,field2,field3"
	got2 := params.Get("$expand")

	if !cmp.Equal(want2, got2) {
		t.Error(cmp.Diff(want2, got2))
	}
}

func TestOptions_SetOrderBy(t *testing.T) {
	opts := QueryOptions{}
	opts.SetOrderBy("field1", OrderByDesc)

	want := "field1 desc"
	got := opts.OrderBy
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	params := opts.ToParams()
	want2 := "field1 desc"
	got2 := params.Get("$orderby")

	if !cmp.Equal(want2, got2) {
		t.Error(cmp.Diff(want2, got2))
	}
}

func TestFilter_SetFilter(t *testing.T) {
	opts := QueryOptions{}
	opts.AddFilter("field1 eq 'something'", "field2 eq 'somethingelse'")

	want := []string{"field1 eq 'something'", "field2 eq 'somethingelse'"}
	got := opts.Filter
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	opts.AddFilter("field3 eq 'another'")
	params := opts.ToParams()

	want2 := "field1 eq 'something' and field2 eq 'somethingelse' and field3 eq 'another'"
	got2 := params.Get("$filter")
	if !cmp.Equal(want2, got2) {
		t.Error(cmp.Diff(want, got2))
	}
}

func TestToParams(t *testing.T) {
	opts := QueryOptions{}
	opts.AddFilter("field1 eq 'something'", "field2 eq 'somethingelse'")
	opts.SetOrderBy("date", OrderByDesc)
	opts.AddSelect("field1", "field2")
	opts.AddExpand("field1", "field2")

	params := opts.ToParams()
	t.Run("Filter", func(t *testing.T) {

		got := params.Get("$filter")
		want := "field1 eq 'something' and field2 eq 'somethingelse'"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("Select", func(t *testing.T) {

		got := params.Get("$select")
		want := "field1,field2"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("Expand", func(t *testing.T) {

		got := params.Get("$expand")
		want := "field1,field2"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("OrderBy", func(t *testing.T) {

		got := params.Get("$orderby")
		want := "date desc"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

}
