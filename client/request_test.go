package dv_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	dv "github.com/erlorenz/dataverse-go/client"
	"github.com/erlorenz/dataverse-go/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestNewRequest_GetBasic(t *testing.T) {

	client, _ := dv.NewClient(dv.Config{
		BaseURL:      "http://example.com",
		TenantID:     "TENANT_ID",
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
	}, dv.WithAuthClient(&mock.AuthClient{}))

	ctx := context.Background()

	req, err := client.NewRequest(ctx, http.MethodGet, "/fake_resources", dv.QueryOptions{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	got := req.Method
	want := http.MethodGet
	if want != got {
		t.Error(cmp.Diff(want, got))
	}

	authHeader := req.Header.Get("Authentication")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Errorf("expected auth header to start with 'Bearer ', got %s", authHeader)
	}

	goturl := strings.Split(req.URL.String(), "?")[0]
	wanturl := "http://example.com/fake_resources"
	if wanturl != goturl {
		t.Error(cmp.Diff(wanturl, goturl))
	}

}

func TestNewRequest_GetWithOptions(t *testing.T) {
	client, _ := dv.NewClient(dv.Config{
		BaseURL:      "http://example.com",
		TenantID:     "TENANT_ID",
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
	}, dv.WithAuthClient(&mock.AuthClient{}))
	opts := dv.QueryOptions{}
	opts.AddExpand("field1")
	opts.AddFilter("field1 eq 'something'")
	opts.AddSelect("field1")
	opts.SetOrderBy("field1", dv.OrderByDesc)

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/fake_resources", opts, nil)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	t.Run("Select", func(t *testing.T) {
		got := req.URL.Query().Get("$select")
		want := "field1"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("Expand", func(t *testing.T) {
		got := req.URL.Query().Get("$expand")
		want := "field1"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("Filter", func(t *testing.T) {
		got := req.URL.Query().Get("$filter")
		want := "field1 eq 'something'"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

	t.Run("OrderBy", func(t *testing.T) {
		got := req.URL.Query().Get("$orderby")
		want := "field1 desc"
		if want != got {
			t.Error(cmp.Diff(want, got))
		}
	})

}

type sometype struct {
	Name string
	Age  int
}

func TestNewRequest_Post(t *testing.T) {
	client, _ := dv.NewClient(dv.Config{
		BaseURL:      "http://example.com",
		TenantID:     "TENANT_ID",
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
	}, dv.WithAuthClient(&mock.AuthClient{}))

	data := sometype{"Fred", 10}

	req, err := client.NewRequest(context.Background(), http.MethodPost, "/fake_resources", dv.QueryOptions{}, data)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}

	defer req.Body.Close()
	body, _ := io.ReadAll(req.Body)

	var got sometype

	json.Unmarshal(body, &got)

	want := data
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
		t.Logf("%#v", body)
	}

}
