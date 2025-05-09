package server

import (
	"context"
	"net"
	"testing"
	"time"
	"vk_test_task/subpub"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"vk_test_task/config"
	pb "vk_test_task/pkg/proto"
)

type mockSubPub struct{}

func (m *mockSubPub) Subscribe(string, subpub.MessageHandler) (subpub.Subscription, error) {
	return &mockSubscription{}, nil
}

func (m *mockSubPub) Publish(string, interface{}) error {
	return nil
}

func (m *mockSubPub) Close(context.Context) error {
	return nil
}

type mockSubscription struct{}

func (m *mockSubscription) Unsubscribe() {}

func TestBasicServerFunctionality(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	defer lis.Close()

	mockSP := &mockSubPub{}
	srv := New(mockSP, config.New())

	err := srv.Start(lis)
	require.NoError(t, err)
	defer srv.Stop(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewPubSubClient(conn)

	_, err = client.Publish(ctx, &pb.PublishRequest{
		Key:  "test",
		Data: "test data",
	})
	require.NoError(t, err)
}
