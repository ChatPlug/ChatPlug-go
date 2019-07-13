package core

type ServiceHandler interface {
	VerifyDepedencies() bool
	LoadService(instanceID string)
	ShutdownService(instanceID string)
}

type ServiceLoader struct {
	App      *App
	handlers map[string]ServiceHandler
}

func (sl *ServiceLoader) Initialize() {
	sl.handlers = make(map[string]ServiceHandler)
	sl.handlers["executable"] = &ExecutableServiceHandler{App: sl.App}
	sl.handlers["node"] = &NodeServiceHandler{App: sl.App}

	for _, v := range sl.handlers {
		v.VerifyDepedencies()
	}
}

func (sl *ServiceLoader) StartupAllInstances() {
	instances := []*ServiceInstance{}
	sl.App.db.Find(&instances)

	for _, instance := range instances {
		instanceService := sl.App.sm.FindServiceWithName(instance.ModuleName)
		sl.handlers[instanceService.Type].LoadService(instance.ID)
		sl.App.sm.LoadInstance(instance.ID)
	}
}

func (sl *ServiceLoader) StartupInstance(instanceID string) {
	var instance ServiceInstance
	sl.App.db.First(&instance, instanceID)

	service := sl.App.sm.FindServiceWithName(instance.ModuleName)

	sl.handlers[service.Type].LoadService(instanceID)
	sl.App.sm.LoadInstance(instance.ID)
}

func (sl *ServiceLoader) ShutdownAllInstances() {
	instances := []*ServiceInstance{}
	sl.App.db.Find(&instances)

	for _, instance := range instances {
		instanceService := sl.App.sm.FindServiceWithName(instance.ModuleName)
		sl.handlers[instanceService.Type].ShutdownService(instance.ID)
	}
}
