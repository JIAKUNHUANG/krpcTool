package main

import (
	"ServiceTest/stub"

	"github.com/JIAKUNHUANG/krpc/server"
)

func main() {
	s := server.CreateService()

	stub.RegisterServiceTestService(s)
	defer s.Listener.Close()

	go s.Service()

	select {}
}
