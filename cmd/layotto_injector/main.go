package main

import (
	log "github.com/sirupsen/logrus"

	"mosn.io/layotto/pkg/injector/service"
)

func main() {
	cfg, err := service.GetConfig()
	if err != nil {
		log.Fatalf("Error getting config: %v", err)
	}
	inj, err := service.NewInjector(cfg)
	if err != nil {
		log.Fatalf("Error creating layotto-injector: %v", err)
	}
	err = inj.Run()
	if err != nil {
		log.Fatalf("Error running layotto-injector: %v", err)
	}
	log.Info("Layotto sidecar injector shutdown gracefully")
}
