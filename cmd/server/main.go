package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"vk_test_task/config"
	"vk_test_task/server"
	"vk_test_task/subpub"
)

func main() {
	cfg := config.New()

	eventBus := subpub.NewSubPub()

	srv := server.New(eventBus, cfg)

	lis, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(cfg.GRPCPort)))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting server on port %s", cfg.GRPCPort)
	if err := srv.Start(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	waitForShutdown(srv, cfg.ShutdownTimeout)
}

func waitForShutdown(srv *server.Server, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
