package dvclient_test

import (
	"errors"
	"testing"

	dataverse "github.com/erlorenz/dataverse-go"
)

func TestNewClient(t *testing.T) {

	config := dataverse.Config{
		BaseURL:      "https://example.crm.com",
		TenantID:     "SOMEID",
		ClientID:     "SOMEID",
		ClientSecret: "somesecret",
	}

	_, err := dataverse.NewClient(config)
	if err != nil {
		t.Errorf("didnt expect error: got %s", err)
	}

}

func TestNewClient_MissingConfig(t *testing.T) {
	baseConfig := dataverse.Config{
		BaseURL:      "https://example.crm.com",
		TenantID:     "SOMEID",
		ClientID:     "SOMEID",
		ClientSecret: "somesecret",
	}

	t.Run("BaseURL", func(t *testing.T) {
		config := baseConfig
		config.BaseURL = ""

		_, err := dataverse.NewClient(config)
		if err == nil {
			t.Errorf("expected error: got nil. %v", config)
		}
	})

	t.Run("TenantID", func(t *testing.T) {
		config := baseConfig
		config.TenantID = ""

		_, err := dataverse.NewClient(config)
		if err == nil {
			t.Errorf("expected error: got nil. %v", config)
		}
	})

	t.Run("ClientID", func(t *testing.T) {
		config := baseConfig
		config.ClientID = ""

		_, err := dataverse.NewClient(config)
		if err == nil {
			t.Errorf("expected error: got nil. %v", config)
		}
	})
	t.Run("ClientSecret", func(t *testing.T) {
		config := baseConfig
		config.ClientSecret = ""

		_, err := dataverse.NewClient(config)
		if err == nil {
			t.Errorf("expected error: got nil. %v", config)
		}
	})

	t.Run("All", func(t *testing.T) {
		config := dataverse.Config{}

		_, err := dataverse.NewClient(config)
		if err == nil {
			t.Errorf("expected error: got nil. %v", config)
		}
		if !errors.Is(err, dataverse.ErrMissingConfig) {
			t.Errorf("expected ErrMissingConfig, got %s", err)
		}
	})

}
