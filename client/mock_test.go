package dataverse_test

import (
	"context"

	dataverse "github.com/erlorenz/dataverse-go/client"
)

// Make sure it satisfies
var _ dataverse.TokenGetter = (*MockAuthClient)(nil)

type MockAuthClient struct {
	ShouldFail bool
}

func (ma *MockAuthClient) GetToken(ctx context.Context) (string, error) {
	if ma.ShouldFail {
		return "", dataverse.ErrTokenAuth
	}
	return "SOME_TOKEN", nil
}
