package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	require := require.New(t)

	repo, err := New()
	require.NoError(err)
	require.NotNil(repo)
	defer repo.Close()

	require.False(repo.Durable())
}
