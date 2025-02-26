package main

import (
	"ServiceFinding/stub"

	"github.com/JIAKUNHUANG/krpc/server"
)

func main() {
	s := server.CreateService()

	s.GetConfig("./config.json")
	stub.RegisterServiceTestService(s)
	defer s.Listener.Close()

	if s.Config.ServiceFindingAdrr.Execute {
		err := s.ServiceFinding()
		if err != nil {
			panic(err)
		}
		defer s.ServiceFindingConn.Close()
	}

	go s.Service()

	select {}
}
