package config_test

import (
	"testing"

	"github.com/mjpitz/highlander-proxy/internal/config"

	"github.com/stretchr/testify/require"
)

func TestRouteSlice(t *testing.T) {
	const goodA = "tcp://0.0.0.0:8080|tcp://localhost:8080"
	const goodB = "tcp://0.0.0.0:8090|tcp://localhost:8090"
	const bad = "invalid:/address/scheme"

	routes := &config.RouteSlice{}

	err := routes.Set(goodA)
	require.Nil(t, err)

	err = routes.Set(goodB)
	require.Nil(t, err)

	err = routes.Set(bad)
	require.NotNil(t, err)
	require.Equal(t, "route: missing pipe", err.Error())

	slice := routes.Routes()
	require.Len(t, slice, 2)

	// entry 1

	require.Equal(t, "tcp", slice[0].From().Scheme)
	require.Equal(t, "0.0.0.0:8080", slice[0].From().Host)
	require.Equal(t, "8080", slice[0].From().Port())

	require.Equal(t, "tcp", slice[0].To().Scheme)
	require.Equal(t, "localhost:8080", slice[0].To().Host)
	require.Equal(t, "8080", slice[0].To().Port())

	// entry 2

	require.Equal(t, "tcp", slice[1].From().Scheme)
	require.Equal(t, "0.0.0.0:8090", slice[1].From().Host)
	require.Equal(t, "8090", slice[1].From().Port())

	require.Equal(t, "tcp", slice[1].To().Scheme)
	require.Equal(t, "localhost:8090", slice[1].To().Host)
	require.Equal(t, "8090", slice[1].To().Port())
}
