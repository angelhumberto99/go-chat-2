package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

func listenServer(msgs *[]string, c net.Conn, name string) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			// terminamos el programa
			os.Exit(1)
		}

		// separación de cadenas
		if strings.Contains(msg, "/file") {
			info := strings.Split(msg, "?")[0]
			data := strings.Split(info, ">")

			fileName := (data[0])[len("/file<"):]
			user := (data[1])[1:len(data[1])-1]
			fileInfo := msg[len(info)+1:]

			if user != name {
				path := "./"+name+"/"+fileName
				dest, err := os.Create(path)
				if err != nil {
					fmt.Println(err)
				} 
				dest.Write([]byte(fileInfo))
				dest.Close()
				*msgs = append(*msgs, user + " envío: " + fileName)
			} else {
				*msgs = append(*msgs, "Enviaste: " + fileName)
			}

		} else if strings.Contains(msg, name+": ") {
			// si el mensaje recibido lleva el nombre del cliente
			// entonces se reemplaza por la palabra "Tú"
			*msgs = append(*msgs, "Tú:" + msg[len(name)+1:])
		} else {
			*msgs = append(*msgs, msg)
		}
	}
}

func listDirectory(name string) string{
	dir, err := ioutil.ReadDir("./"+name)
	var files []string
	input := bufio.NewScanner(os.Stdin)

    if err != nil {
        fmt.Println(err)
    }
	for _, v := range dir {
		files = append(files, v.Name())
	}

	fmt.Println("Archivos")
	if len(files) == 0 {
		fmt.Println("Usted no cuenta con ningún archivo")
		return "-1"
	}
	// menú
	for i, v := range files {
		fmt.Printf("%d) %s\n", i, v)
	}
	fmt.Printf("%d) Salir\n", len(files))
	input.Scan()
	resp,_ := strconv.Atoi(input.Text())
	if len(files) == resp {
		return "-1"
	}
	return files[resp]
}

func sendFile(fileName, userName string, c net.Conn) {
	origin, err := os.ReadFile("./"+userName+"/"+fileName)
	if err != nil {
		fmt.Println(err)
	}
	// información para enviar
	fileStr := "/file<"+fileName+">("+userName+")?"+string(origin) 
	err = gob.NewEncoder(c).Encode(fileStr)
	if err != nil {
		fmt.Println(err)
	}
}

func client(port string, chat string) {
	var msgs []string
	var name string
	menu := chat + "\n" +
			"1) Mostrar mensajes/archivos\n" + 
			"2) Enviar mensaje\n" +
			"3) Enviar archivo\n" + 
			"4) Salir\n"
	input := bufio.NewScanner(os.Stdin)
	
	
	// pedimos el nombre al usuario
	fmt.Print("Ingrese su nombre: ")
	input.Scan()
	name = input.Text()

	// creamos una carpeta para sus archivos
	os.Mkdir(name, os.ModePerm)

	// conectamos el cliente al servidor
	c, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
	}

	// escuchara todas las respuestas del servidor
	go listenServer(&msgs, c, name)

	for {
		fmt.Print(menu)
		input.Scan()

		switch input.Text() {
			case "1": // mostrar mensajes
				fmt.Println("Mensajes")
				for _,v := range msgs {
					fmt.Println(v)
				}
			case "2": // enviar mensaje
				msg := bufio.NewScanner(os.Stdin)
				fmt.Print(">>> ")
				msg.Scan()
				err = gob.NewEncoder(c).Encode(name +": "+ msg.Text())
				if err != nil {
					fmt.Println(err)
				}
			case "3": // enviar archivo
				fileName := listDirectory(name)
				if fileName != "-1" {
					sendFile(fileName, name, c)
				}
			case "4": // terminar cliente
				fmt.Println("Terminando cliente")

			default:
				fmt.Println("Opción incorrecta")
		}

		if input.Text() == "4" {
			// se envia un mensaje al servidor 
			// para eliminar la conexion al cliente
			err = gob.NewEncoder(c).Encode("/quit")
			if err != nil {
				fmt.Println(err)
			}
			break
		}

	}
	// terminamos la llamada
	c.Close()
}

func getOptions() []string {
	c, err := rpc.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}

	var res []string
	aux := ""
	err = c.Call("Server.GetChats", aux, &res)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func getPort(option string) string {
	c, err := rpc.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}

	var res string
	err = c.Call("Server.GetPort", option, &res)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	options := getOptions()

	// enlistamos las salas de chat por RPC
	for i, v := range(options) {
		fmt.Printf("%d) %s\n", i+1, v)
	}
	fmt.Println("Elige una sala de chat: ")
	scanner.Scan()
	opc,_ := strconv.Atoi(scanner.Text())

	// obtenemos el puerto de la sala de chat por RPC
	port := getPort(options[opc-1])

	client(port, options[opc-1])
}