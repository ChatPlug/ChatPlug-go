package core

import "context"

func (r *Resolver) Thread() ThreadResolver {
	return &threadResolver{r}
}

type threadResolver struct{ *Resolver }

func (r *threadResolver) Service(ctx context.Context, obj *Thread) (*ServiceInstance, error) {
	var instance ServiceInstance

	r.App.db.First(&instance, "id = ?", obj.ServiceInstanceID)
	return &instance, nil
}
