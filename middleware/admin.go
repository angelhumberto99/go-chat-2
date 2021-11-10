package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func main() {
	c, err := net.Dial("tcp", ":9998")
	if err != nil {
		fmt.Println(err)
	}
	err = gob.NewEncoder(c).Encode("/admin")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			// terminamos el programa
			os.Exit(1)
		}
		fmt.Println(msg)
	}
}