package main

import (
	"../model"
	"./grep"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Server Server Struct
type Server struct {
	Protocol, Port string
}

// HelloWorld Hello World RPC exapmple
func (s *Server) HelloWorld(args *model.RPCArgs, reply *string) error {
	*reply = grep.Grep(args.A)
	return nil
}

// This function will register and initiate server
func main() {
	server := new(Server)
	server.Port = "8080"
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+server.Port)
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	http.Serve(l, nil)
}
