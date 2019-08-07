package core

import "context"

func (r *Resolver) ServiceInstance() ServiceInstanceResolver {
	return &serviceInstanceResolver{r}
}

type serviceInstanceResolver struct{ *Resolver }

func (r *serviceInstanceResolver) Service(ctx context.Context, obj *ServiceInstance) (*Service, error) {
	var service Service

	r.App.db.First(&service, "name = ?", obj.ModuleName)
	return &service, nil
}
