package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/erlorenz/dataverse-go/internal/auth"
)

func TestGetToken(t *testing.T) {

	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	authClient, err := auth.New(tenantID, clientID, clientSecret)
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
