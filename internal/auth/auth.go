package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

// ErrTokenAuth is returned when the AuthClient cannot get a token.
var ErrTokenAuth = errors.New("cannot retrieve auth token")

// MSALClient is used to retrieve an AccessToken.
// Implements the TokenGetter interface.
type MSALClient struct {
	cca    confidential.Client
	scopes []string
	Logger *slog.Logger
}

// New creates an [MSALClient] using a client secret.
func New(tenantID string, clientID string, clientSecret string) (*MSALClient, error) {
	if tenantID == "" || clientID == "" || clientSecret == "" {
		missing := []string{}
		switch {
		case tenantID == "":
			missing = append(missing, "tenantID")
		case clientID == "":
			missing = append(missing, "clientID")
		case clientSecret == "":
			missing = append(missing, "clientSecret")
		}
		return nil, fmt.Errorf("missing params %s", strings.Join(missing, ", "))
	}

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		err = fmt.Errorf("dv credential: %w", err)
		return nil, err
	}

	authority := "https://login.microsoft.com/" + tenantID

	confidentialClient, err := confidential.New(authority, clientID, cred)
	if err != nil {
		err = fmt.Errorf("new confidentialClient error: %w", err)
		return nil, err
	}

	// Don't think there is any reason to use a different one.
	// Can have this as a config param if ever need to.
	scopes := []string{"https://api.businesscentral.dynamics.com/.default"}

	return &MSALClient{
		cca:    confidentialClient,
		scopes: scopes,
		Logger: slog.Default(),
	}, nil

}

// GetToken implements TokenGetter. It returns an access token.
func (ac *MSALClient) GetToken(ctx context.Context) (string, error) {

	ac.Logger.Debug("Acquiring token...")
	result, err := ac.cca.AcquireTokenSilent(ctx, ac.scopes)
	if err != nil {
		// cache miss, authenticate with another AcquireToken... method
		ac.Logger.Debug("Cache miss, calling AcquireTokenByCredential...")

		result, err = ac.cca.AcquireTokenByCredential(ctx, ac.scopes)
		if err != nil {
			return "", fmt.Errorf("error getting access token: %w", err)
		}
	}
	ac.Logger.Debug("Successfully acquired token.")
	return result.AccessToken, nil
}

// Debug sets the logger to log at slog.LevelDebug.
// If w is nil it defaults to os.Stdout.
func (ac *MSALClient) Debug(w io.Writer) {
	var out io.Writer = os.Stdout
	if w != nil {
		out = w
	}

	ac.Logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
