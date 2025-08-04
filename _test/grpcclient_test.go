package _test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/grpcclient"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func newBufConnListener() *bufconn.Listener {
	return bufconn.Listen(bufSize)
}

func dialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestClient_GetConn(t *testing.T) {
	listener := newBufConnListener()

	// Dummy Server starten
	s := grpc.NewServer()
	go func() {
		if err := s.Serve(listener); err != nil {

			t.Fatalf("Server failed: %v", err)
		}
	}()
	t.Cleanup(func() {
		s.Stop()
	})

	factory := grpcclient.NewClient(
		grpc.WithContextDialer(dialer(listener)),
		grpc.WithInsecure(), // nur für Tests
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	)

	conn1, err := factory.GetConnection("dummy:1234")
	require.NoError(t, err)
	require.NotNil(t, conn1)

	conn2, err := factory.GetConnection("dummy:1234")
	require.NoError(t, err)
	require.Same(t, conn1, conn2, "cached connection should be reused")
}
