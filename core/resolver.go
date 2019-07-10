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

func (r *mutationResolver) SendMessage(ctx context.Context, input NewMessage) (*Message, error) {
	panic("not implemented")
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

func (r *mutationResolver) AddThreadToGroup(ctx context.Context, input *NewThread) (*ThreadGroup, error) {
	panic("not implemented")
}
func (r *mutationResolver) SetInstanceStatus(ctx context.Context, instanceID string, status *InstanceStatus) (*ServiceInstance, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateNewInstance(ctx context.Context, serviceModuleName string, instanceName string) (*ServiceInstance, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Messages(ctx context.Context) ([]*Message, error) {
	panic("not implemented")
}
func (r *queryResolver) Instances(ctx context.Context) ([]*ServiceInstance, error) {
	panic("not implemented")
}

func (r *queryResolver) ThreadGroups(ctx context.Context) ([]*ThreadGroup, error) {
	groups := []*ThreadGroup{}
	r.App.db.Find(&groups)
	return groups, nil
}

func (r *queryResolver) Services(ctx context.Context) ([]*Service, error) {
	return r.App.sm.services, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessageReceived(ctx context.Context, threadID string) (<-chan *Message, error) {
	panic("not implemented")
}
