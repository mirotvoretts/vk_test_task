package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vk_test_task/pkg/proto"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubSubClient(conn)

	testSubscribe(client)
	testPublish(client)
}

func testSubscribe(client pb.PubSubClient) {
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{Key: "test"})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	go func() {
		for {
			event, err := stream.Recv()
			if err != nil {
				log.Printf("Stream ended: %v", err)
				return
			}
			log.Printf("Received event: %s", event.Data)
		}
	}()
}

func testPublish(client pb.PubSubClient) {
	_, err := client.Publish(context.Background(), &pb.PublishRequest{
		Key:  "test",
		Data: "Hello from updated client",
	})
	if err != nil {
		log.Fatalf("Publish failed: %v", err)
	}
	log.Println("Message published successfully")

	time.Sleep(2 * time.Second)
}
