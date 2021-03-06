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

	var handler ServiceHandler
	switch service.Type {
	case "node":
		handler = &NodeServiceHandler{BaseServiceHandler: baseHandler}
	case "python3":
		handler = &Python3ServiceHandler{BaseServiceHandler: baseHandler}
	case "executable":
		handler = &ExecutableServiceHandler{BaseServiceHandler: baseHandler}
	default:
		log.Printf("Error: Service type %s unknown!", service.Type)
		// os.Exit(1)
	}

	if err := handler.Init(); err != nil {
		log.Printf(err.Error())
		// os.Exit(1)
	}

	sl.handlers[service.Name] = handler

	return handler
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
		sl.StartupInstance(instance.ID)
	}
}

func (sl *ServiceLoader) GetHandlerForInstance(instanceID string) ServiceHandler {
	var instance ServiceInstance
	sl.App.db.Where("id = ?", instanceID).First(&instance)

	service := sl.App.sm.FindServiceWithName(instance.ModuleName)
	if service == nil {
		log.Printf("Service %s not found!", instance.ModuleName)
		return nil
	}
	sl.App.sm.LoadInstance(instance.ID)
	return sl.GetHandlerForService(service)
}

func (sl *ServiceLoader) StartupInstance(instanceID string) {
	sl.GetHandlerForInstance(instanceID).LoadService(instanceID)
}

func (sl *ServiceLoader) ShutdownAllInstances() {
	instances := []*ServiceInstance{}
	sl.App.db.Find(&instances)

	for _, instance := range instances {
		service := sl.App.sm.FindServiceWithName(instance.ModuleName)
		sl.GetHandlerForService(service).ShutdownService(instance.ID)
	}
}
