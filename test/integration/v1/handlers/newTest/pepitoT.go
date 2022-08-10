package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func NonTest(t *testing.T) {
	t.Run("test pepito", func(t *testing.T) {
		err := fmt.Errorf("error test")
		require.Error(t, err)
	})
}
