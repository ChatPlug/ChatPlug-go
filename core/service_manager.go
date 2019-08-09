package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type LoadedInstance struct {
	instanceID       string
	messageEventBroadcaster *EventBroadcaster
	searchRequestEventBroadcaster *EventBroadcaster
	searchResponseEventBroadcaster *EventBroadcaster
}

type ServiceManager struct {
	services  []*Service
	instances []*LoadedInstance
}

func (sm *ServiceManager) LoadAvailableServices() {
	files, err := ioutil.ReadDir("services")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		manifestPath := "services/" + f.Name() + "/manifest.json"
		if f.IsDir() && Exists(manifestPath) {
			jsonFile, err := os.Open(manifestPath)
			if err == nil {
				byteVal, _ := ioutil.ReadAll(jsonFile)
				var service Service
				json.Unmarshal(byteVal, &service)

				log.Printf("Service found! %s", service.DisplayName)
				sm.services = append(sm.services, &service)
			}
		}
	}
}

func (sm *ServiceManager) FindServiceWithName(moduleName string) *Service {
	for _, n := range sm.services {
		if n.Name == moduleName {
			return n
		}
	}
	return nil
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (sm *ServiceManager) LoadInstance(instanceID string) *LoadedInstance {
	if !sm.IsInstanceLoaded(instanceID) {
		instance := &LoadedInstance{
			instanceID:       instanceID,
			messageEventBroadcaster: NewEventBroadcaster(),
			searchRequestEventBroadcaster: NewEventBroadcaster(),
			searchResponseEventBroadcaster: NewEventBroadcaster(),
		}

		sm.instances = append(sm.instances, instance)
		return instance
	}

	return nil
}

func (sm *ServiceManager) FindLoadedInstance(instanceID string) *LoadedInstance {
	for _, n := range sm.instances {
		if n.instanceID == instanceID {
			return n
		}
	}
	return nil
}

func (sm *ServiceManager) IsInstanceLoaded(instanceID string) bool {
	for _, n := range sm.instances {
		if n.instanceID == instanceID {
			return true
		}
	}
	return false
}
