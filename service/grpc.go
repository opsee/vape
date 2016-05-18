package service

import (
	opsee "github.com/opsee/basic/service"
	"google.golang.org/grpc"
)

type service struct {
	Server *grpc.Server
}

func New() *service {
	s := &service{}

	server := grpc.NewServer()
	opsee.RegisterVapeServer(server, s)

	s.Server = server

	return s
}
