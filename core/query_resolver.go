package core

import "context"

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
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
