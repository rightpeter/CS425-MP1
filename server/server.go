package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"

	"CS425/CS425-MP1/model"
	"CS425/CS425-MP1/server/grep"
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

func (s *Server) setIP(IP string) {
	s.config.Current.IP = IP
}

func (s *Server) getPort() int {
	return s.config.Current.Port
}

func (s *Server) setPort(port int) {
	s.config.Current.Port = port
}

func (s *Server) getFilePath() string {
	return s.config.Current.LogPath
}

func (s *Server) setFilePath(path string) {
	s.config.Current.LogPath = path
}

// Grep RPC to call grep on server
func (s *Server) Grep(args *model.RPCArgs, reply *string) (err error) {
	*reply, err = grep.Grep(args.Commands, s.getFilePath())
	return err
}

// This function will register and initiate server
func main() {

	// parse argument
	configFilePath := flag.String("c", "./config.json", "Config file path")
	port := flag.Int("p", 8080, "Port number")
	IP := flag.String("ip", "0.0.0.0", "IP address")

	flag.Parse()

	// load config file
	configFile, e := ioutil.ReadFile(*configFilePath)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	// Class for rpc
	server := newServer()
	server.loadConfigFromJSON(configFile)

	server.setIP(*IP)
	server.setPort(*port)
	server.setFilePath(server.getFilePath())

	fmt.Printf("Starting server on IP: %s and port: %d", *IP, *port)

	// init the rpc server
	newServer := rpc.NewServer()
	newServer.Register(server)

	l, e := net.Listen("tcp", fmt.Sprintf(":%d", server.getPort()))
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	// server start
	newServer.Accept(l)
}
