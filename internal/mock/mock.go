package mock

import (
	"context"

	"github.com/erlorenz/dataverse-go/internal/auth"
)

type AuthClient struct {
	ShouldFail bool
}

func (ma *AuthClient) GetToken(ctx context.Context) (string, error) {
	if ma.ShouldFail {
		return "", auth.ErrTokenAuth
	}
	return "SOME_TOKEN", nil
}
