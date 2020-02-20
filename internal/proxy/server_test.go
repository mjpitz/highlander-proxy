package proxy_test

import (
	"github.com/mjpitz/highlander-proxy/internal/election"
	"github.com/mjpitz/highlander-proxy/internal/proxy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer(t *testing.T) {
	leader := election.NewLeader()

	server := &proxy.Server{
		Protocol: "tcp",
		BindAddress: "localhost:8080",
		RemoteAddress: "localhost:8090",
		Leader: leader,
	}

	a, ctxa, err := server.GetForwardAddress()
	require.NotNil(t, err)
	require.Equal(t, "no leader elected", err.Error())
	require.Nil(t, ctxa)
	require.Equal(t, "", a)

	// when I'm not leader, make sure I point to the leader

	leader.Update("localhost:1234")

	b, ctxb, err := server.GetForwardAddress()
	require.Nil(t, err)
	require.NotNil(t, ctxb)
	require.Equal(t, "localhost:1234", b)

	// when I'm leader, make sure I point to my remote

	leader.Update("localhost:8080")

	c, ctxc, err := server.GetForwardAddress()
	require.Nil(t, err)
	require.NotNil(t, ctxc)
	require.Equal(t, "localhost:8090", c)

	select {
	case <-ctxb.Done():
		break
	default:
		require.Fail(t, "ctxb did not terminate")
	}
}
