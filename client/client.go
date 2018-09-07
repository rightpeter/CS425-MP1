package client

import (
	"MP0/server"
	"fmt"
	"log"
	"net/rpc"
)

// Client struct
type Client struct {
	ServerAddress, PortNumber string
	client                    *rpc.Client
}

var serverAddress = "localhost"

// CallRPC example
func (c *Client) CallRPC(str string) error {
	args := &server.Args{A: str}

	var reply string
	err := c.client.Call("Server.HelloWorld", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Resp: ", reply)
	return nil
}

// RegisterClient example
func (c *Client) RegisterClient() error {
	client, err := rpc.DialHTTP("tcp", c.ServerAddress+":"+c.PortNumber)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c.client = client
	return err
}
