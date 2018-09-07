package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Server Server Struct
type Server struct {
	Protocol, Port string
}

// Args Input Arguments for rpc
type Args struct {
	A string
}

// HelloWorld Hello World RPC exapmple
func (s *Server) HelloWorld(args *Args, reply *string) error {
	*reply = args.A
	return nil
}

// RegisterServer This function will register and initiate server
func (s *Server) RegisterServer() {
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+s.Port)
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	go http.Serve(l, nil)
}
