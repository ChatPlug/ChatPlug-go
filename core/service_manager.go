package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ServiceManager struct {
	services []*Service
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

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
