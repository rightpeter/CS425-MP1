package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"CS425/CS425-MP1/model"

	"encoding/json"
)

// Client struct
type Client struct {
	config model.NodesConfig
}

func (c *Client) loadConfigFromJSON(jsonFile []byte) error {
	return json.Unmarshal(jsonFile, &c.config)
}

func (c *Client) callRPC(client *rpc.Client, commands []string) (string, error) {
	args := model.RPCArgs{Commands: commands}
	var reply string
	err := client.Call("Server.Grep", &args, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (c *Client) distribuitedGrep(clientID int, commands []string) model.RPCResult {
	result := model.RPCResult{ClientID: clientID}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.config.Nodes[clientID].IP, c.config.Nodes[clientID].Port), 300*time.Millisecond)
	if err != nil {
		result.Alive = false
		return result
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	result.Alive = true

	reply, err := c.callRPC(client, commands)
	if err != nil {
		result.Error = err
		return result
	}
	result.Reply = reply
	return result

}

// DistributedGrep non blocking distributed grep
func (c *Client) DistributedGrep(commands []string) {
	ch := make(chan model.RPCResult)
	for k := range c.config.Nodes {
		go func(clientID int) {
			select {
			case ch <- c.distribuitedGrep(clientID, commands):
			case <-time.After(time.Second):
				ch <- model.RPCResult{ClientID: clientID, Alive: false}
			}
		}(k)
	}

	for range c.config.Nodes {
		result := <-ch
		fmt.Println(strings.Repeat("+", 30) + "[VM" + strconv.Itoa(result.ClientID) + "]" + strings.Repeat("+", 30))

		if !result.Alive {
			fmt.Printf("VM%d died!\n\n", result.ClientID)
		} else if result.Error != nil {
			if result.Error.Error() == "exit status 1" {
				fmt.Printf("Lines Count: 0\n\n")
			} else {
				fmt.Printf("Grep fail! Error: %v\n\n", result.Error)
			}
		} else {
			r := strings.NewReader(result.Reply)
			scanner := bufio.NewScanner(r)

			line := 0
			for scanner.Scan() {
				fmt.Println(scanner.Text())
				line++
			}

			fmt.Printf("Lines Count: %d\n\n", line)
		}
	}
}

func main() {
	configFile, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	c := &Client{}
	c.loadConfigFromJSON(configFile)

	c.DistributedGrep(os.Args[1:])
}
