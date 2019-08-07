package core

import (
	"context"
	"fmt"
)

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessageReceived(ctx context.Context) (<-chan *MessagePayload, error) {
	instance := r.App.InstanceForContext(ctx)
	if instance == nil {
		return nil, fmt.Errorf("Access denied")
	}

	eventBroadcaster := r.App.sm.FindEventBoardcasterByInstanceID(instance.ID)

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

func (r *subscriptionResolver) ConfigurationReceived(ctx context.Context, configurationRequest ConfigurationRequest) (<-chan *ConfigurationResponse, error) {
	configurationRequest.resChan = make(chan *ConfigurationResponse)
	r.App.ch.configurationQueue <- &configurationRequest

	return configurationRequest.resChan, nil
}
