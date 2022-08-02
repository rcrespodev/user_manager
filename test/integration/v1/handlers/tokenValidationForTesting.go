package handlers

import (
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TokenValidationForTesting(t *testing.T, header http.Header) {
	token := header.Get("Token")
	require.NotNil(t, token)
	_, err := kernel.Instance.Jwt().ValidateToken(token)
	require.NoError(t, err)
}
