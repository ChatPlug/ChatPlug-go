package core

import (
	"context"
	"log"
)

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SendMessage(ctx context.Context, instanceID string, input MessageInput) (*Message, error) {
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
			avatarURL := input.Author.AvatarURL
			if avatarURL == "" {
				avatarURL = "https://i.imgur.com/3yPh9fE.png"
			}

			r.App.db.Where("origin_id = ?", input.OriginID).FirstOrCreate(&messageAuthor,
				MessageAuthor{
					Username:  input.Author.Username,
					OriginID:  input.Author.OriginID,
					AvatarURL: avatarURL,
				})

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

func (r *mutationResolver) DeleteServiceInstance(ctx context.Context, id string) (string, error) {
	r.App.db.Where("id = ?", id).Delete(&ServiceInstance{})
	return id, nil
}

func (r *mutationResolver) DeleteThread(ctx context.Context, id string) (string, error) {
	r.App.db.Where("id = ?", id).Delete(&Thread{})
	return id, nil
}

func (r *mutationResolver) AddThreadToGroup(ctx context.Context, input *ThreadInput) (*ThreadGroup, error) {
	var group ThreadGroup

	r.App.db.First(&group, "id = ?", input.GroupID)

	r.App.db.Model(&group).Association("Threads").Append(&Thread{
		OriginID:          input.OriginID,
		Name:              input.Name,
		ServiceInstanceID: input.InstanceID,
	})

	return &group, nil
}

func (r *mutationResolver) SetInstanceStatus(ctx context.Context, instanceID string, status *InstanceStatus) (*ServiceInstance, error) {
	var instance ServiceInstance

	r.App.db.First(&instance, "id = ?", instanceID)

	log.Printf("Instance %s status changed: %s -> %s", instanceID, instance.Status, status.String())
	if instance.Status == "STOPPED" && status.String() == "RUNNING" {
		r.App.sl.StartupInstance(instanceID)
	}
	if instance.Status == "RUNNING" && status.String() == "STOPPED" {
		r.App.sl.GetHandlerForInstance(instanceID).ShutdownService(instanceID)
	}

	r.App.db.Model(&instance).Update("Status", status.String())

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
