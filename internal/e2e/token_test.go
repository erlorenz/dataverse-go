package e2e

import (
	"context"
	"testing"

	dataverse "github.com/erlorenz/dataverse-go/client"
)

func TestGetToken(t *testing.T) {
	authClient, err := dataverse.NewSecretAuthClient(TenantID, ClientID, ClientSecret)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
		t.Log(authClient)
	}

	token, err := authClient.GetToken(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	if token == "" {
		t.Error("token is empty")
	}
}
