package core

import (
	"sync"
)

//EventBroadcaster is used to broadcast data to many listeners on separate goroutines. Used to exchange messages in thread groups
type EventBroadcaster struct {
	subscribers      []chan interface{}
	subscribersMutex *sync.RWMutex
}

func NewEventBroadcaster() *EventBroadcaster {
	return &EventBroadcaster{
		subscribers:      []chan interface{}{},
		subscribersMutex: &sync.RWMutex{},
	}
}

func (eb *EventBroadcaster) Subscribe() (eventChannel chan interface{}, cancel func()) {
	eventChannel = make(chan interface{})

	eb.subscribersMutex.Lock()
	defer eb.subscribersMutex.Unlock()
	eb.subscribers = append(eb.subscribers, eventChannel)
	didCancel := false
	cancel = func() {
		if didCancel {
			panic("subscription already cancelled")
		}
		didCancel = true
		eb.subscribersMutex.Lock()
		defer eb.subscribersMutex.Unlock()
		b := eb.subscribers[:0]
		// remove the channel from the subscribers slice
		for _, sub := range eb.subscribers {
			if sub != eventChannel {
				b = append(b, sub)
			}
		}
		eb.subscribers = b
		close(eventChannel)
	}
	return
}

func (eb *EventBroadcaster) Broadcast(ev interface{}) {
	eb.subscribersMutex.RLock()
	defer eb.subscribersMutex.RUnlock()
	for _, sub := range eb.subscribers {
		sub <- ev
	}
}
