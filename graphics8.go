//
//  (C) Copyright 2021  Pavel Tisnovsky
//
//  All rights reserved. This program and the accompanying materials
//  are made available under the terms of the Eclipse Public License v1.0
//  which accompanies this distribution, and is available at
//  http://www.eclipse.org/legal/epl-v10.html
//
//  Contributors:
//      Pavel Tisnovsky
//

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handler(connection net.Conn) {
	for {
		fmt.Fprintf(connection, "?\n")
		status, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Print("No response!")
			return
		} else {
			log.Print(status)
		}
	}
}

func server(listener net.Listener) {
	log.Print("Waiting for connections")
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Print("Connection refused!")
		}
		defer connection.Close()
		go handler(connection)
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Can't open the port!", err)
	}
	defer l.Close()
	server(l)
}
