package e2e

import (
	"context"
	"testing"

	"github.com/erlorenz/dataverse-go"
)

type account struct {
	ID string `json:"accountid"`
}

func (a account) Validate() error {
	return nil
}

func TestClient(t *testing.T) {
	config := getConfig(t)

	client, err := dataverse.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	qo := dataverse.QueryOptions{Top: 3}
	qo.AddSelect("accountid")

	req, err := client.NewRequest(context.Background(), "GET", "accounts", qo, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("header: %s", req.Header.Get("Authorization")[:20]+"...")
	t.Logf("url: %s", req.URL.String())
	t.Logf("method: %s", req.Method)
	t.Logf("query: %s", req.URL.RawQuery)

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	decoded, err := dataverse.DecodeResponseList[account](res)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(decoded)
	if len(decoded) != 3 {
		t.Fatalf("expected len of 3: %#v", decoded)
	}
}
