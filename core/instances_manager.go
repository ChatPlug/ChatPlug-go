package core

type LoadedInstance struct {
	instanceID       string
	eventBroadcaster *EventBroadcaster
}

type ServiceInstancesManager struct {
	instances []*LoadedInstance
}

func (sim *ServiceInstancesManager) LoadInstance(instanceID string) *LoadedInstance {
	if !sim.IsInstanceLoaded(instanceID) {
		instance := &LoadedInstance{
			instanceID:       instanceID,
			eventBroadcaster: NewEventBroadcaster(),
		}

		sim.instances = append(sim.instances, instance)
		return instance
	}

	return nil
}

func (sim *ServiceInstancesManager) FindEventBoardcasterByInstanceID(instanceID string) *EventBroadcaster {
	for _, n := range sim.instances {
		if n.instanceID == instanceID {
			return n.eventBroadcaster
		}
	}
	return nil
}

func (sim *ServiceInstancesManager) IsInstanceLoaded(instanceID string) bool {
	for _, n := range sim.instances {
		if n.instanceID == instanceID {
			return true
		}
	}
	return false
}
