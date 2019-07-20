package core

import (
	"context"
	"log"
)

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

func (r *subscriptionResolver) ConfigurationReceived(ctx context.Context, configurationRequest ConfigurationRequest) (<-chan *ConfigurationResponse, error) {
	log.Println("A O S")
	configurationRequest.resChan = make(chan *ConfigurationResponse)
	r.App.ch.configurationQueue <- &configurationRequest

	return configurationRequest.resChan, nil
}
