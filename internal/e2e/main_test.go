package e2e

import (
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var (
	TenantID     string
	ClientID     string
	ClientSecret string
	BaseURL      string
)

func TestMain(m *testing.M) {

	godotenv.Load("./.env")

	TenantID = os.Getenv("TENANT_ID")
	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	BaseURL = os.Getenv("BASE_URL")

	if TenantID == "" || ClientID == "" || ClientSecret == "" || BaseURL == "" {
		slog.Error("Missing config!")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
