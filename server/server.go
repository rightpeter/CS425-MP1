package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"../model"
	"./grep"
)

// Server Server Struct
type Server struct {
	config model.NodesConfig
}

func newServer() *Server {
	return &Server{}
}

func (s *Server) loadConfigFromJSON(jsonFile []byte) error {
	return json.Unmarshal(jsonFile, &s.config)
}

func (s *Server) getIP() string {
	return s.config.Current.IP
}

func (s *Server) getPort() int {
	return s.config.Current.Port
}

// Grep RPC to call grep on server
func (s *Server) Grep(args *model.RPCArgs, reply *string) error {
	*reply = grep.Grep(args.Command)
	return nil
}

// This function will register and initiate server
func main() {
	configFile, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	server := newServer()
	server.loadConfigFromJSON(configFile)

	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", fmt.Sprintf(":%d", server.getPort()))
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	http.Serve(l, nil)
}
