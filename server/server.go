package server

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"vk_test_task/config"
	pb "vk_test_task/pkg/proto"
	"vk_test_task/subpub"
)

type Server struct {
	pb.UnimplementedPubSubServer
	subpub     subpub.SubPub
	config     *config.Config
	grpcServer *grpc.Server
	wg         sync.WaitGroup
}

func New(sp subpub.SubPub, cfg *config.Config) *Server {
	return &Server{
		subpub: sp,
		config: cfg,
	}
}

func (s *Server) Start(lis net.Listener) error {
	s.grpcServer = grpc.NewServer()
	pb.RegisterPubSubServer(s.grpcServer, s)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.grpcServer.Serve(lis); err != nil {
			return
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	}
}

func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	msgChan := make(chan string, 100)

	sub, err := s.subpub.Subscribe(req.Key, func(msg interface{}) {
		if data, ok := msg.(string); ok {
			select {
			case msgChan <- data:
			default:
			}
		}
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case data := <-msgChan:
			if err := stream.Send(&pb.Event{Data: data}); err != nil {
				return status.Errorf(codes.Internal, "failed to send event: %v", err)
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *Server) Publish(_ context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	if err := s.subpub.Publish(req.Key, req.Data); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish: %v", err)
	}
	return &emptypb.Empty{}, nil
}
