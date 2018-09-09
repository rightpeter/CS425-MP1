package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"

	"../model"
	"encoding/json"
)

// Client struct
type Client struct {
	clients map[int]*rpc.Client
	config  model.NodesConfig
}

func newClient() *Client {
	return &Client{clients: make(map[int]*rpc.Client)}
}

func (c *Client) loadConfigFromJSON(jsonFile []byte) error {
	return json.Unmarshal(jsonFile, &c.config)
}

func (c *Client) registerClient() (err error) {
	for _, v := range c.config.Nodes {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", v.IP, v.Port))
		if err != nil {
			log.Println("dialing: ", err)
			break
		}
		c.clients[v.ID] = client
	}
	return err
}

func (c *Client) callRPC(serverID int, args interface{}, reply interface{}) error {
	err := c.clients[serverID].Call("Server.HelloWorld", args, reply)
	return err
}

func main() {
	configFile, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	c := newClient()
	c.loadConfigFromJSON(configFile)

	err := c.registerClient()
	if err != nil {
		log.Fatal("error registering client:", err)
	}

	args := &model.RPCArgs{A: "say this!!!"}
	var reply string
	err = c.callRPC(1, args, &reply)
	if err != nil {
		log.Fatal("Call RPC Failed: ", err)
	}

	log.Println("Call RPC Suceed: ", reply)
}
