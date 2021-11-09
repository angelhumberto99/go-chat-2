package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Server struct {}

func (s *Server) GetChats(args string, reply *[]string) error {
	args = ""
	*reply = []string{"politica", "educación", "jocosidades"}
	return nil
}

func (s *Server) GetPort(chat string, reply *string) error {
	chats := make(map[string]string)
	chats["politica"] = ":9000"
	chats["educación"] = ":9001"
	chats["jocosidades"] = ":9002"
	*reply = chats[chat]
	return nil
}

func server() {
	// nombre del chat: usuarios conectados
	// chats := make(map[string]int)

	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()
	
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	var input string
	fmt.Scanln(&input)
}