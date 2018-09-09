package main

import (
	"CS425-MP1/model"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"

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
	err := c.clients[serverID].Call("Server.Grep", args, reply)
	return err
}

func (c *Client) distributedGrep(args interface{}, reply interface{}) error {
	for _, v := range c.clients {
		err := v.Call("Server.Grep", args, reply)
		if err != nil {
			log.Fatal("Error calling Server.Grep: ", err)
		}
	}
	return nil
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

	// Take input from user

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		fmt.Print("> ")
		input := scanner.Text()
		fmt.Println(input)
		args := &model.RPCArgs{Command: input}
		var reply string
		err := c.distributedGrep(args, &reply)
		if err != nil {
			log.Fatal("Call RPC Failed: ", err)
		}
		log.Println("Call RPC Suceed: ", reply)
	}

}
