package subpub

import (
	"context"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type subscription struct {
	subject    string
	handler    MessageHandler
	msgChan    chan interface{}
	cancel     context.CancelFunc
	mu         sync.Mutex
	closed     bool
	removeFunc func(*subscription)
}

type subPub struct {
	mu          sync.RWMutex
	subscribers map[string][]*subscription
	closed      bool
	wg          sync.WaitGroup
}

func NewSubPub() SubPub {
	return &subPub{
		subscribers: make(map[string][]*subscription),
	}
}

func (s *subPub) Subscribe(subject string, handler MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, context.Canceled
	}

	ctx, cancel := context.WithCancel(context.Background())
	sub := &subscription{
		subject:    subject,
		handler:    handler,
		msgChan:    make(chan interface{}, 100),
		cancel:     cancel,
		removeFunc: s.removeSub,
	}

	s.subscribers[subject] = append(s.subscribers[subject], sub)
	s.wg.Add(1)

	go s.processMessages(ctx, sub)
	return sub, nil
}

func (s *subPub) processMessages(ctx context.Context, sub *subscription) {
	defer s.wg.Done()
	for {
		select {
		case msg := <-sub.msgChan:
			sub.handler(msg)
		case <-ctx.Done():
			return
		}
	}
}

func (s *subPub) removeSub(sub *subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := s.subscribers[sub.subject]
	for i, cur := range subs {
		if cur == sub {
			s.subscribers[sub.subject] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

func (s *subPub) Publish(subject string, msg interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return context.Canceled
	}

	for _, sub := range s.subscribers[subject] {
		go func(s *subscription) {
			select {
			case s.msgChan <- msg:
			default:
			}
		}(sub)
	}

	return nil
}

func (s *subPub) Close(ctx context.Context) error {
	s.mu.Lock()
	s.closed = true
	subs := s.subscribers
	s.subscribers = nil
	s.mu.Unlock()

	for _, subjectSubs := range subs {
		for _, sub := range subjectSubs {
			sub.Unsubscribe()
		}
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *subscription) Unsubscribe() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	s.cancel()
	s.removeFunc(s)
	close(s.msgChan)
}
