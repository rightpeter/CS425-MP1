package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"sort"
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

func (c *Client) distribuitedGrep(clientConfig model.NodeConfig, commands []string) model.RPCResult {
	result := model.RPCResult{ClientID: clientConfig.ID}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", clientConfig.IP, clientConfig.Port), 300*time.Millisecond)
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
	for _, v := range c.config.Nodes {
		// Use goroutine for non-blocking grep
		go func(clientConfig model.NodeConfig) {
			select {
			case ch <- c.distribuitedGrep(clientConfig, commands):
			case <-time.After(500 * time.Second):
				ch <- model.RPCResult{ClientID: clientConfig.ID, Alive: false}
			}
		}(v)
	}

	summary := make(map[int]int)
	// collect results from pipeline ch, and print it out
	for range c.config.Nodes {
		result := <-ch
		fmt.Println(strings.Repeat("+", 30) + "[VM" + strconv.Itoa(result.ClientID) + "]" + strings.Repeat("+", 30))

		if !result.Alive {
			fmt.Printf("VM%d died!\n\n", result.ClientID)
			summary[result.ClientID] = -2
		} else if result.Error != nil {
			if result.Error.Error() == "exit status 1" {
				fmt.Printf("Lines Count: 0\n\n")
				summary[result.ClientID] = 0
			} else {
				fmt.Printf("Grep fail! Error: %v\n\n", result.Error)
				summary[result.ClientID] = -1
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
			summary[result.ClientID] = line
		}
	}
	// sort keys to print summary from VM0 to VM10
	var keys []int
	for k := range summary {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	total := 0
	for _, k := range keys {
		if summary[k] == -2 {
			fmt.Printf("VM%d: Failed, ", k)
		} else if summary[k] == -1 {
			fmt.Printf("VM%d: Grep Error, ", k)
		} else {
			total += summary[k]
			fmt.Printf("VM%d: %d, ", k, summary[k])
		}
	}
	fmt.Printf("\nTotal line count: %d\n", total)
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
