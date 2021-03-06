package proxy_test

import (
	"context"
	"io"
	"testing"

	"github.com/mjpitz/highlander-proxy/internal/proxy"

	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	message := "Hello World!"

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	readerA, writerA := io.Pipe()
	readerB, writerB := io.Pipe()

	go proxy.Pipe(ctx, readerA, writerB)

	nA, err := writerA.Write([]byte(message))
	require.Nil(t, err)

	data := make([]byte, nA)
	nB, err := readerB.Read(data)
	require.Nil(t, err)

	require.Equal(t, nA, nB)
	require.Equal(t, message, string(data))
}
