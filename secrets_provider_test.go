package observability_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/observability"
	"github.com/xaionaro-go/secret"
	"golang.org/x/oauth2"
)

func TestParseSecretsFrom(t *testing.T) {
	type someStruct struct {
		String string
		Token  *secret.Any[oauth2.Token]
	}
	sample := &someStruct{
		String: "1",
		Token: ptr(secret.New(oauth2.Token{
			AccessToken:  "5",
			TokenType:    "6",
			RefreshToken: "7",
		})),
	}
	secrets := observability.ParseSecretsFrom(sample)
	require.Equal(t, []string{"5", "6", "7"}, secrets)
}
