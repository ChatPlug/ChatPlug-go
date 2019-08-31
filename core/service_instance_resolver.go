package core

import "context"

func (r *Resolver) ServiceInstance() ServiceInstanceResolver {
	return &serviceInstanceResolver{r}
}

type serviceInstanceResolver struct{ *Resolver }

func (r *serviceInstanceResolver) Service(ctx context.Context, obj *ServiceInstance) (*Service, error) {
	return r.App.sm.FindServiceWithName(obj.ModuleName), nil
}
