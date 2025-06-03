package e2e

import (
	"context"
	"testing"

	"github.com/erlorenz/dataverse-go/internal/auth"
)

func TestGetToken(t *testing.T) {
	authClient, err := auth.New(TenantID, ClientID, ClientSecret)
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
