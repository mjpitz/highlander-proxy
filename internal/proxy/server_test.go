package proxy_test

import (
	"testing"

	"github.com/mjpitz/highlander-proxy/internal/config"
	"github.com/mjpitz/highlander-proxy/internal/election"
	"github.com/mjpitz/highlander-proxy/internal/proxy"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	leader := election.NewLeader()

	route := &config.Route{}
	err := route.Set("tcp://localhost:8080|tcp://localhost:8090")
	require.Nil(t, err)

	server := &proxy.Server{
		Route:    route,
		Identity: "8.0.8.0",
		Leader:   leader,
	}

	aproto, a, ctxa, err := server.GetForwardAddress()
	require.NotNil(t, err)
	require.Equal(t, "no leader elected", err.Error())
	require.Nil(t, ctxa)
	require.Equal(t, "", aproto)
	require.Equal(t, "", a)

	// when I'm not leader, make sure I point to the leader

	leader.Update("1.2.3.4")

	bproto, b, ctxb, err := server.GetForwardAddress()
	require.Nil(t, err)
	require.NotNil(t, ctxb)
	require.Equal(t, "tcp", bproto)
	require.Equal(t, "1.2.3.4:8080", b)

	// when I'm leader, make sure I point to my remote

	leader.Update("8.0.8.0")

	cproto, c, ctxc, err := server.GetForwardAddress()
	require.Nil(t, err)
	require.NotNil(t, ctxc)
	require.Equal(t, "tcp", cproto)
	require.Equal(t, "localhost:8090", c)

	select {
	case <-ctxb.Done():
		break
	default:
		require.Fail(t, "ctxb did not terminate")
	}
}
