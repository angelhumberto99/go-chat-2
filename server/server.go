package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

func listenClients(msgs *[]string, clients *[]net.Conn, c net.Conn, channel chan string) {
	var msg string
	for {
		// recibimos los mensajes del cliente
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		// si llega "/quit" eliminamos la conexion del cliente
		if msg == "/quit" {
			for i, v := range *clients {
				if v == c {
					*clients = append((*clients)[:i], (*clients)[i+1:]...)
					break
				}
			}
			c.Close()
			return
		}

		if strings.Contains(msg, "/file") {
			info := strings.Split(msg, "?")[0]
			data := strings.Split(info, ">")

			fileName := (data[0])[len("/file<"):]
			user := (data[1])[1:len(data[1])-1]
			*msgs = append(*msgs, user + " envío: "+ fileName)
		} else {
			// añadimos el mensaje al slice de mensajes
			*msgs = append(*msgs, msg)
		}
		// enviamos el mensaje a todos los clientes
		channel <- msg
	}
}

func checkConnection(s net.Listener, clients *[]net.Conn, msgs *[]string, channel chan string) {
	var exists bool
	for {
		// peticiones del cliente
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
		}
		
		exists = false
		for _,v := range(*clients) {
			if v == c {
				exists = true
				break
			}
		}
		
		// si el cliente no existe, entonces lo añadimos
		if !exists {
			*clients = append(*clients, c)
			go listenClients(msgs, clients, c, channel)
		}
	}
}

func clientsMsgsHandler(clients *[]net.Conn, channel chan string) {
	for {
		// si recibimos una señal por el canal
		// enviaremos a todos los clientes un mensaje
		msg := <-channel
		for _,c := range *clients {
			err := gob.NewEncoder(c).Encode(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func saveMsgs(msgs []string) {
	file, err := os.Create("dataBase.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for _, v := range msgs {
		file.WriteString(v+"\n")
	}
}

// func sendInfoMidware(port, chat string, clients *[]net.Conn) {
// 	c, err := net.Dial("tcp", ":9999")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func main() {
	var clients []net.Conn
	var msgs []string
	channel := make(chan string)
	menu := "1) Mostrar mensajes/archivos\n" + 
			"2) Respaldar mensajes\n" + 
			"3) Salir\n"
	input := bufio.NewScanner(os.Stdin)

	fmt.Println("Ingrese el puerto del servidor: ")
	input.Scan()

	port := ":" + input.Text()

	fmt.Println("Ingrese la tematica del chat: ")
	input.Scan()

	// go sendInfoMidware(port, input.Text(), &clients)

	// se crea el servidor
	s, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()
	go checkConnection(s, &clients, &msgs, channel)
	go clientsMsgsHandler(&clients, channel)

	for {
		fmt.Print(menu)
		input.Scan()
		switch input.Text() {
			case "1": // mostrar mensajes
				fmt.Println("Mensajes")
				for _,v := range msgs {
					fmt.Println(v)
				}
			case "2": // respaldar mensajes
				saveMsgs(msgs)
			case "3": // terminar cliente
				fmt.Println("Terminando Servidor")
			default:
				fmt.Println("Opción incorrecta")
		}
		if input.Text() == "3" {
			break
		}
	}
}