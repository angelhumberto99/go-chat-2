package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type Server struct {}

func (s *Server) GetChats(args string, reply *[]string) error {
	args = ""
	*reply = []string{"Politica", "Deportes", "Economia"}
	return nil
}

func (s *Server) GetPort(chat string, reply *string) error {
	chats := make(map[string]string)
	chats["Politica"] = ":9001"
	chats["Deportes"] = ":9002"
	chats["Economia"] = ":9003"
	*reply = chats[chat]
	return nil
}

func listenServer(c net.Conn) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			// terminamos el programa
			os.Exit(1)
		}
		fmt.Println("server info: ", msg)
	}
}

func checkServers(servers *[]net.Conn) {
	var exists bool
	s, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println(err)
	}
	defer s.Close()
	for {
		// peticiones del cliente
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
		}
		
		exists = false
		for _,v := range(*servers) {
			if v == c {
				exists = true
				break
			}
		}
		
		// si el cliente no existe, entonces lo añadimos
		if !exists {
			*servers = append(*servers, c)
			go listenServer(c)
		}
	}
}

func rpcServer() {
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
	var servers []net.Conn
	go rpcServer()
	go checkServers(&servers)
	var input string
	fmt.Scanln(&input)
}