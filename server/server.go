package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
func (s *Server) Grep(args *model.RPCArgs, reply *string) error {
	*reply = grep.Grep(args.Command, s.getFilePath())
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

	port := flag.Int("p", server.getPort(), "Port number")
	IP := flag.String("ip", server.getIP(), "IP address")
	logPath := flag.String("f", server.getFilePath(), "Log file path")

	flag.Parse()

	fmt.Printf("Starting server on IP: %s and port: %d", *IP, *port)

	server.setIP(*IP)
	server.setPort(*port)
	server.setFilePath(*logPath)

	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", fmt.Sprintf(":%d", server.getPort()))
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	http.Serve(l, nil)
}
