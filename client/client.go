package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"

	"CS425/CS425-MP1/model"

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

func (c *Client) callRPC(serverID int, command string, chReply chan<- string, chErr chan<- error) {
	args := &model.RPCArgs{Command: command}
	var reply string
	err := c.clients[serverID].Call("Server.Grep", args, &reply)
	fmt.Println("Calling grep!")
	if err != nil {
		chErr <- err
		return
	}
	fmt.Println("returning grep!")
	chReply <- reply
}

func (c *Client) distributedGrep(command string) string {
	replies := make(chan string)
	errors := make(chan error)
	var reply string

	for k := range c.clients {
		go c.callRPC(k, command, replies, errors)
	}

	// append replies
	for i := 0; i < len(c.clients); i++ {
		select {
		case rep := <-replies:
			reply += rep
		case err := <-errors:
			log.Println("Error grepping: ", err)
		}
	}
	return reply
}

func main() {

	fmt.Println("Starting client...")

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
		input := scanner.Text()
		fmt.Println(input)
		reply := c.distributedGrep(input)
		log.Println("Call RPC Suceeded: ", reply)
		fmt.Print("> ")
	}

}
