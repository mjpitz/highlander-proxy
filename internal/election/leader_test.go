package election_test

import (
	"testing"

	"github.com/mjpitz/highlander-proxy/internal/election"

	"github.com/stretchr/testify/require"
)

func TestLeader(t *testing.T) {
	leader := election.NewLeader()
	require.NotNil(t, leader)

	a, ctxa, ok := leader.Get()
	require.Equal(t, false, ok)
	require.Nil(t, ctxa)
	require.Equal(t, "", a)

	leader.Update("b")

	b, ctxb, ok := leader.Get()
	require.Equal(t, true, ok)
	require.NotNil(t, ctxb)
	require.Equal(t, "b", b)

	leader.Update("c")

	c, ctxc, ok := leader.Get()
	require.Equal(t, true, ok)
	require.NotNil(t, ctxc)
	require.Equal(t, "c", c)

	select {
	case <-ctxb.Done():
		break
	default:
		require.Fail(t, "ctxb did not terminate")
	}
}
