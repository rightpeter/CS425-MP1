package main

import (
	"MP0/client"
	"MP0/server"
	"log"
)

var serverAddress = "localhost"

func main() {

	server := new(server.Server)
	server.Port = "1234"
	server.RegisterServer()

	c := new(client.Client)
	c.PortNumber = "1234"
	c.ServerAddress = "localhost"
	err := c.RegisterClient()

	if err != nil {
		log.Fatal("error registering client:", err)
	}
	c.CallRPC("say this!!!")

	// // Server
	// arith := new(server.Arith)
	// rpc.Register(arith)
	// rpc.HandleHTTP()

	// l, e := net.Listen("tcp", ":1234")
	// if e != nil {
	// 	log.Fatal("listen error: ", e)
	// }

	// go http.Serve(l, nil)

	// // Server

	// // Client
	// client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	// if err != nil {
	// 	log.Fatal("dialing:", err)
	// }

	// args := &server.Args{7, 9}

	// var reply int
	// err = client.Call("Arith.Multiply", args, &reply)
	// if err != nil {
	// 	log.Fatal("arith error:", err)
	// }
	// fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
	// // Client
}
