package core

import (
	"context"
	"log"
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
	var lastSentMsg *Message
	r.App.db.Preload("Threads").Find(&groups)

	for _, group := range groups {
		rightGroup := false
		var rightThreadID string
		for _, thread := range group.Threads {

			if thread.ServiceInstanceID == instanceID && thread.OriginID == input.OriginThreadID {
				rightGroup = true
				rightThreadID = thread.ID
			}
		}

		if rightGroup {
			log.Println(group.Name)
			var messageAuthor MessageAuthor
			r.App.db.Where("origin_id = ?", input.OriginID).FirstOrCreate(&messageAuthor, MessageAuthor{Username: input.Author.Username, OriginID: input.Author.OriginID})

			attachments := make([]Attachment, 0)

			for _, attachment := range input.Attachments {
				attachmentObj := Attachment{
					Type:      attachment.Type,
					OriginID:  attachment.OriginID,
					SourceURL: attachment.SourceURL,
				}

				attachments = append(attachments, attachmentObj)
			}

			msg := &Message{
				OriginID:        input.OriginID,
				Body:            input.Body,
				ThreadID:        rightThreadID,
				MessageAuthorID: messageAuthor.ID,
				Attachments:     attachments,
			}

			r.App.db.Model(&group).Association("Messages").Append(msg)

			for _, thread := range group.Threads {
				if msg.ThreadID != thread.ID {
					r.App.sm.FindEventBoardcasterByInstanceID(thread.ServiceInstanceID).Broadcast(&MessagePayload{
						TargetThreadID: thread.OriginID,
						Message:        msg,
					})
				}

			}
			lastSentMsg = msg
		}
	}
	return lastSentMsg, nil
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

	if status.String() == "RUNNING" && instance.Status != "RUNNING" {
		r.App.sl.StartupInstance(instanceID)
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
	r.App.db.Preload("Attachments").Find(&messages)
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

func (r *Resolver) Message() MessageResolver {
	return &messageResolver{r}
}

type messageResolver struct{ *Resolver }

func (r *messageResolver) Author(ctx context.Context, obj *Message) (*MessageAuthor, error) {
	var author MessageAuthor

	r.App.db.First(&author, "id = ?", obj.MessageAuthorID)
	return &author, nil
}

func (r *messageResolver) Thread(ctx context.Context, obj *Message) (*Thread, error) {
	var thread Thread

	r.App.db.First(&thread, "id = ?", obj.ThreadID)
	return &thread, nil
}
