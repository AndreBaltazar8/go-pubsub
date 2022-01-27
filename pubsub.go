package pubsub

import "sync"

type Topic[K any] interface {
	Subscribe() <-chan K
	Unsubscribe(<-chan K)
	Publish(value K)
	Close()
}

type topicImpl[K any] struct {
	lock        sync.RWMutex
	subscribers []chan K
	closed      bool
}

func NewTopic[K any]() Topic[K] {
	return &topicImpl[K]{
		subscribers: make([]chan K, 0),
	}
}

func (topic *topicImpl[K]) Subscribe() <-chan K {
	topic.lock.Lock()
	defer topic.lock.Unlock()

	ch := make(chan K, 1)
	topic.subscribers = append(topic.subscribers, ch)
	return ch
}

func (topic *topicImpl[K]) Unsubscribe(ch <-chan K) {
	topic.lock.Lock()
	defer topic.lock.Unlock()

	for i, c := range topic.subscribers {
		if c == ch {
			topic.subscribers = append(topic.subscribers[:i], topic.subscribers[i+1:]...)
			close(c)
			return
		}
	}
}

func (topic *topicImpl[K]) Publish(value K) {
	topic.lock.RLock()
	defer topic.lock.RUnlock()

	if topic.closed {
		return
	}

	for _, ch := range topic.subscribers {
		ch <- value
	}
}

func (topic *topicImpl[K]) Close() {
	topic.lock.Lock()
	defer topic.lock.Unlock()

	if !topic.closed {
		topic.closed = true
		for _, ch := range topic.subscribers {
			close(ch)
		}
	}
}
