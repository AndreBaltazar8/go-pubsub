package pubsub

import (
	"sync"
	"testing"
)

func TestTopic(t *testing.T) {
	topic := NewTopic[string]()
	topicChannel := topic.Subscribe()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		received := <-topicChannel
		if received != "hello" {
			t.Errorf("Expected 'hello', got '%s'", received)
		}
		wg.Done()
	}()
	topic.Publish("hello")
	wg.Wait()
}
