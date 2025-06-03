package e2e

import (
	"log"
	"os"
	"testing"

	"github.com/erlorenz/dataverse-go"
)

func getConfig(t *testing.T) dataverse.Config {
	t.Helper()

	var cfg dataverse.Config

	cfg.TenantID = os.Getenv("TENANT_ID")
	cfg.ClientID = os.Getenv("CLIENT_ID")
	cfg.ClientSecret = os.Getenv("CLIENT_SECRET")
	cfg.BaseURL = os.Getenv("BASE_URL")

	if cfg.TenantID == "" || cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.BaseURL == "" {
		log.Fatal("Missing config!")
	}

	return cfg
}
