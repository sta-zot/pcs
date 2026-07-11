package main

import (
	"fmt"
	"github/sta-zot/pcs/internal/config"
	"github/sta-zot/pcs/internal/provisioning"
	"github/sta-zot/pcs/internal/registrator"
)

type svc struct {
}

func (s *svc) GetConfig(mac, ip string, info ...string) ([]byte, error) {
	return nil, nil
}

func main() {
	//root := chi.NewRouter()

	fmt.Printf("SIPServer - %s\n", config.GlobalConfig.Sip.GetServer())
	router := registrator.New("/api/config")
	s := &svc{}
	_ = provisioning.NewHandler(s, router)
	// http.ListenAndServe(":8080", router.Handler())
	// fmt.Println("Server listen on port :8080")
}
