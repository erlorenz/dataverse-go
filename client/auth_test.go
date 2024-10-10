package dv_test

import (
	"testing"

	dv "github.com/erlorenz/dataverse-go/client"
)

var (
	fakeTenantID     = "FAKE_TENANT_ID"
	fakeClientID     = "FAKE_CLIENT_ID"
	fakeClientSecret = "FAKE_CLIENT_SECRET"
)

func TestNewAuth(t *testing.T) {
	_, err := dv.NewSecretAuthClient(fakeTenantID, fakeClientID, fakeClientSecret)
	if err != nil {
		t.Error(err)
	}
}

func TestNewAuth_MissingInfo(t *testing.T) {
	t.Parallel()

	t.Run("TenantID", func(t *testing.T) {
		_, err := dv.NewSecretAuthClient("", fakeClientID, fakeClientSecret)
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ClientID", func(t *testing.T) {
		_, err := dv.NewSecretAuthClient(fakeTenantID, "", fakeClientSecret)
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ClientID", func(t *testing.T) {
		_, err := dv.NewSecretAuthClient(fakeTenantID, fakeClientID, "")
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
