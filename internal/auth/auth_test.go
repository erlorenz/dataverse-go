package auth_test

import (
	"testing"

	"github.com/erlorenz/dataverse-go/internal/auth"
)

var (
	fakeTenantID     = "FAKE_TENANT_ID"
	fakeClientID     = "FAKE_CLIENT_ID"
	fakeClientSecret = "FAKE_CLIENT_SECRET"
)

func TestNewAuth(t *testing.T) {
	_, err := auth.New(fakeTenantID, fakeClientID, fakeClientSecret)
	if err != nil {
		t.Error(err)
	}
}

func TestNewAuth_MissingInfo(t *testing.T) {
	t.Parallel()

	t.Run("TenantID", func(t *testing.T) {
		_, err := auth.New("", fakeClientID, fakeClientSecret)
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ClientID", func(t *testing.T) {
		_, err := auth.New(fakeTenantID, "", fakeClientSecret)
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("ClientID", func(t *testing.T) {
		_, err := auth.New(fakeTenantID, fakeClientID, "")
		t.Log(err)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
