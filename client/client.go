package main

import (
	"../model"
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
	args := &model.RPCArgs{A: str}

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

func main() {
	c := new(Client)
	c.PortNumber = "8080"
	c.ServerAddress = "localhost"
	err := c.RegisterClient()

	if err != nil {
		log.Fatal("error registering client:", err)
	}
	c.CallRPC("say this!!!")
}
