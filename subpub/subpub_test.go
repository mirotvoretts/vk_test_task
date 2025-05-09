package subpub

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSubscribePublish(t *testing.T) {
	sp := NewSubPub()
	var wg sync.WaitGroup
	wg.Add(1)

	sub, err := sp.Subscribe("test", func(msg interface{}) {
		defer wg.Done()
		if msg != "hello" {
			t.Errorf("Expected 'hello', got %v", msg)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	sp.Publish("test", "hello")
	wg.Wait()
}

func TestClose(t *testing.T) {
	sp := NewSubPub()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	sp.Subscribe("test", func(msg interface{}) {
		defer wg.Done()
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		sp.Close(ctx)
	}()

	sp.Publish("test", "data")
	wg.Wait()
}
