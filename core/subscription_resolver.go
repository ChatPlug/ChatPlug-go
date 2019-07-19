package core

import "context"

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessageReceived(ctx context.Context, instanceID string) (<-chan *MessagePayload, error) {
	eventBroadcaster := r.App.sm.FindEventBoardcasterByInstanceID(instanceID)

	messages := make(chan *MessagePayload, 1)
	go func() {
		msgChan, cancel := eventBroadcaster.Subscribe()
	Loop:
		for {
			select {
			case msg := <-msgChan:
				if msg != nil {
					messages <- msg
				}
			case <-ctx.Done():
				cancel()
				break Loop
			}
		}
	}()
	return messages, nil
}