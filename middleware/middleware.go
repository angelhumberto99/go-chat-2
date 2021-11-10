package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Server struct {}

// tematica: cantidad de clientes conectados
var clients map[string]int

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

func (s *Server) GetClients(chat string, reply *int) error {
	*reply = clients[chat]
	return nil
}

func listenServer(c net.Conn, servers *[]net.Conn, infoChan chan string) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if strings.Contains(msg, "/quit") {
			for i, v := range *servers {
				if v == c {
					*servers = append((*servers)[:i], (*servers)[i+1:]...)
					break
				}
			}
			c.Close()
			infoChan <- "Offline: "+ msg[len("/quit"):]
			return
		}
		if err != nil {
			fmt.Println(err)
			// terminamos el programa
			os.Exit(1)
		}
		
		topic := strings.Split(msg, " ->")[0]
		online := strings.Split(msg, "(")[1]
		online = strings.Split(online, " clientes)")[0]

		val,_ := strconv.Atoi(online)
		clients[topic] = val
		infoChan <- "Online: "+ msg
	}
}

func listenAdmin(c net.Conn, infoChan chan string) {
	for {
		val := <-infoChan
		err := gob.NewEncoder(c).Encode(val)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func checkServers(servers *[]net.Conn) {
	var exists bool
	infoChan := make(chan string)
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
			var msg string
			_ = gob.NewDecoder(c).Decode(&msg)
			if msg == "/admin" {
				fmt.Println("admin")
				go listenAdmin(c, infoChan)
			} else {
				*servers = append(*servers, c)
				fmt.Println("servidor")
				go listenServer(c, servers, infoChan)
				infoChan <- "Online: "+ msg
			}
		}
	}
}

func rpcServer() {
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
	clients = make(map[string]int)
	go rpcServer()
	go checkServers(&servers)
	var input string
	fmt.Scanln(&input)
}