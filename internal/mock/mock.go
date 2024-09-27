package mock

import (
	"context"

	dataverse "github.com/erlorenz/dataverse-go/client"
)

type AuthClient struct {
	ShouldFail bool
}

func (ma *AuthClient) GetToken(ctx context.Context) (string, error) {
	if ma.ShouldFail {
		return "", dataverse.ErrTokenAuth
	}
	return "SOME_TOKEN", nil
}
