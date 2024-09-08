package main

import (
	"fmt"
	"log"
	"net"
	file "test_tages/proto"
	"test_tages/server/constants"
	"test_tages/server/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", constants.ServerPort)
	if err != nil {
		log.Fatalf(err.Error())
	}

	s := grpc.NewServer()
	fileService := service.NewFileService(constants.DataDir)
	file.RegisterFileServiceServer(s, fileService)

	fmt.Printf("Server started, listening on %s \n", constants.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
