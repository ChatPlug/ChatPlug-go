package core

import (
	"log"
)

type ServiceLoader struct {
	App      *App
	handlers map[string]ServiceHandler
}

func (sl *ServiceLoader) Initialize() {
	sl.handlers = make(map[string]ServiceHandler)
}

func (sl *ServiceLoader) CreateServiceHandler(service *Service) ServiceHandler {
	baseHandler := &BaseServiceHandler{App: sl.App, Service: service}

	switch service.Type {
	case "node":
		sl.handlers[service.Name] = &NodeServiceHandler{BaseServiceHandler: baseHandler}
	case "executable":
		sl.handlers[service.Name] = &ExecutableServiceHandler{BaseServiceHandler: baseHandler}
	default:
		log.Printf("Error: Service type %s unknown!", service.Type)
		// crash?
	}

	sl.handlers[service.Name].VerifyDepedencies()
	// TODO: crash when not verified?

	return sl.handlers[service.Name]
}

func (sl *ServiceLoader) GetHandlerForService(service *Service) ServiceHandler {
	if handler, ok := sl.handlers[service.Name]; ok {
		return handler
	}

	return sl.CreateServiceHandler(service)
}

func (sl *ServiceLoader) StartupAllInstances() {
	instances := []*ServiceInstance{}
	sl.App.db.Find(&instances)

	for _, instance := range instances {
		service := sl.App.sm.FindServiceWithName(instance.ModuleName)
		sl.GetHandlerForService(service).LoadService(instance.ID)

		sl.App.sm.LoadInstance(instance.ID)
	}
}

func (sl *ServiceLoader) StartupInstance(instanceID string) {
	var instance ServiceInstance
	sl.App.db.Where("id = ?", instanceID).First(&instance)

	service := sl.App.sm.FindServiceWithName(instance.ModuleName)
	sl.GetHandlerForService(service).LoadService(instanceID)

	sl.App.sm.LoadInstance(instance.ID)
}

func (sl *ServiceLoader) ShutdownAllInstances() {
	instances := []*ServiceInstance{}
	sl.App.db.Find(&instances)

	for _, instance := range instances {
		service := sl.App.sm.FindServiceWithName(instance.ModuleName)
		sl.GetHandlerForService(service).ShutdownService(instance.ID)
	}
}
