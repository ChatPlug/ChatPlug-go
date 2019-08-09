package core

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
)

// Mutation resolver handles all graphql mutations
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SetSearchResponse(ctx context.Context, forQuery string, threads []*ThreadSearchResultInput) (*SearchResponse, error) {
	instance := r.App.InstanceForContext(ctx)
	if instance == nil {
		return &SearchResponse{}, fmt.Errorf("Access denied")
	}

	loadedInstance := r.App.sm.FindLoadedInstance(instance.ID)

	threadResults := make([]*ThreadSearchResult, 0)
	for _, thread := range threads {
		threadResults = append(threadResults, &ThreadSearchResult{
			Name:     thread.Name,
			OriginID: thread.OriginID,
			IconURL:  thread.IconURL,
		})
	}

	res := &SearchResponse{ForQuery: forQuery, Threads: threadResults}

	loadedInstance.searchResponseEventBroadcaster.Broadcast(res)

	return res, nil
}

func (r *mutationResolver) SearchThreadsInService(ctx context.Context, q string, instanceID string) (*SearchResponse, error) {
	loadedInstance := r.App.sm.FindLoadedInstance(instanceID)

	loadedInstance.searchRequestEventBroadcaster.Broadcast(&SearchRequest{Query: q})

	var res *SearchResponse
	msgChan, cancel := loadedInstance.searchResponseEventBroadcaster.Subscribe()
	for msg := range msgChan {
		if msg.(*SearchResponse).ForQuery == q {
			cancel()
			res = msg.(*SearchResponse)
			break
		}
	}

	return res, nil
}

func (r *mutationResolver) SendMessage(ctx context.Context, input MessageInput) (*Message, error) {
	instance := r.App.InstanceForContext(ctx)
	if instance == nil {
		return &Message{}, fmt.Errorf("Access denied")
	}

	var groups []ThreadGroup
	var lastSentMsg *Message
	r.App.db.Preload("Threads").Find(&groups)

	for _, group := range groups {
		rightGroup := false
		var rightThreadID string
		for _, thread := range group.Threads {

			if thread.ServiceInstanceID == instance.ID && thread.OriginID == input.OriginThreadID {
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
				if thread.Readonly != nil && *thread.Readonly == true {
					continue
				}
				if msg.ThreadID != thread.ID {
					r.App.sm.FindLoadedInstance(thread.ServiceInstanceID).messageEventBroadcaster.Broadcast(&MessagePayload{
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

	iconURL := "https://i.imgur.com/3yPh9fE.png"
	if input.IconURL != nil {
		iconURL = *input.IconURL
	}

	r.App.db.Model(&group).Association("Threads").Append(&Thread{
		OriginID:          input.OriginID,
		Name:              input.Name,
		ServiceInstanceID: input.InstanceID,
		Readonly:          input.Readonly,
		IconURL:           iconURL,
	})

	return &group, nil
}

func (r *mutationResolver) SetInstanceStatus(ctx context.Context, status *InstanceStatus) (*ServiceInstance, error) {
	instance := r.App.InstanceForContext(ctx)
	if instance == nil {
		return &ServiceInstance{}, fmt.Errorf("Access denied")
	}

	log.Printf("Instance %s status changed: %s -> %s", instance.ID, instance.Status, status.String())
	if instance.Status == "STOPPED" && status.String() == "RUNNING" {
		r.App.sl.StartupInstance(instance.ID)
	}
	if instance.Status == "RUNNING" && status.String() == "STOPPED" {
		r.App.sl.GetHandlerForInstance(instance.ID).ShutdownService(instance.ID)
	}

	r.App.db.Model(&instance).Update("Status", status.String())

	return instance, nil
}

func (r *mutationResolver) CreateNewInstance(ctx context.Context, serviceModuleName string, instanceName string) (*NewServiceInstanceCreated, error) {
	token := GenerateAccessToken()
	newInstance := &ServiceInstance{
		Name:        instanceName,
		ModuleName:  serviceModuleName,
		Threads:     []Thread{},
		AccessToken: token,
	}

	r.App.db.Create(newInstance)
	res := &NewServiceInstanceCreated{
		Instance:    newInstance,
		AccessToken: token,
	}
	return res, nil
}

// GenerateAccessToken generates a random token used to authenticate services
func GenerateAccessToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
