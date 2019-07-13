package core

import (
	"context"
)

type Resolver struct {
	App *App
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SendMessage(ctx context.Context, instanceID string, input NewMessage) (*Message, error) {
	var groups []ThreadGroup
	r.App.db.Preload("Threads").Find(&groups)

	for _, group := range groups {
		rightGroup := false
		var rightThreadID string
		for _, thread := range group.Threads {
			if thread.ServiceInstanceID == instanceID {
				rightGroup = true
				rightThreadID = thread.ID
			}
		}

		if rightGroup {
			var messageAuthor MessageAuthor
			r.App.db.Where("origin_id = ?", input.OriginID).FirstOrCreate(&messageAuthor, MessageAuthor{Username: input.Author.Username, OriginID: input.Author.OriginID})
			msg := &Message{
				OriginID:        input.OriginID,
				Body:            input.Body,
				ThreadID:        rightThreadID,
				MessageAuthorID: messageAuthor.ID,
			}

			r.App.db.Model(&group).Association("Messages").Append(msg)

			for _, thread := range group.Threads {
				if thread.ServiceInstanceID != instanceID {
					r.App.sim.FindEventBoardcasterByInstanceID(thread.ServiceInstanceID).Broadcast(&MessagePayload{
						TargetThreadID: thread.ID,
						Message:        msg,
					})
				}
			}
			return msg, nil
		}
	}
	return nil, nil
}

func (r *mutationResolver) CreateThreadGroup(ctx context.Context, name string) (*ThreadGroup, error) {
	group := &ThreadGroup{
		Name:     name,
		Messages: []Message{},
		Threads:  []Thread{},
	}

	r.App.db.Create(group)
	return group, nil
}

func (r *mutationResolver) DeleteThreadGroup(ctx context.Context, id string) (string, error) {
	r.App.db.Where("id = ?", id).Delete(&ThreadGroup{})
	return id, nil
}

func (r *mutationResolver) DeleteThread(ctx context.Context, id string) (string, error) {
	r.App.db.Where("id = ?", id).Delete(&Thread{})
	return id, nil
}

func (r *mutationResolver) AddThreadToGroup(ctx context.Context, input *NewThread) (*ThreadGroup, error) {
	var group ThreadGroup

	r.App.db.First(&group, "id = ?", input.GroupID)

	r.App.db.Model(&group).Association("Threads").Append(&Thread{
		OriginID:          input.OriginID,
		Name:              input.Name,
		ServiceInstanceID: input.ServiceID,
	})

	return &group, nil
}

func (r *mutationResolver) SetInstanceStatus(ctx context.Context, instanceID string, status *InstanceStatus) (*ServiceInstance, error) {
	var instance ServiceInstance

	r.App.db.First(&instance, "id = ?", instanceID)
	r.App.db.Model(&instance).Update("Status", status.String())

	if status.String() == "RUNNING" {
		r.App.sm.StartupServiceInstance(&instance)
		r.App.sim.LoadInstance(instance.ID)
	}

	return &instance, nil
}

func (r *mutationResolver) CreateNewInstance(ctx context.Context, serviceModuleName string, instanceName string) (*ServiceInstance, error) {
	newInstance := &ServiceInstance{
		Name:       instanceName,
		ModuleName: serviceModuleName,
		Threads:    []Thread{},
	}

	r.App.db.Create(newInstance)
	return newInstance, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Messages(ctx context.Context) ([]*Message, error) {
	messages := []*Message{}
	r.App.db.Find(&messages)
	return messages, nil
}

func (r *queryResolver) Instances(ctx context.Context) ([]*ServiceInstance, error) {
	instances := []*ServiceInstance{}
	r.App.db.Find(&instances)
	return instances, nil
}

func (r *queryResolver) ThreadGroups(ctx context.Context) ([]*ThreadGroup, error) {
	groups := []*ThreadGroup{}
	r.App.db.Preload("Threads").Find(&groups)
	return groups, nil
}

func (r *queryResolver) Services(ctx context.Context) ([]*Service, error) {
	return r.App.sm.services, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessageReceived(ctx context.Context, instanceID string) (<-chan *MessagePayload, error) {
	eventBroadcaster := r.App.sim.FindEventBoardcasterByInstanceID(instanceID)
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
