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

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 640
	height = 480
)

func handler(connection net.Conn, command chan string, quit chan bool, listener net.Listener) {
	log.Print("Listener handler")
	for {
		fmt.Fprintf(connection, "?\n")
		status, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Print("No response!")
			return
		} else {
			log.Print(status)
			if status == "quit\n" {
				log.Print("Closing handler")
				listener.Close()
				quit <- true
				command <- status
				return
			}
		}
	}
}

func server(listener net.Listener, command chan string) {
	log.Print("Waiting for connections")
	quit := make(chan bool)

	for {
		select {
		case <-quit:
			println("quitting")
			return
		default:
		}
		log.Print("Connection listener")
		connection, err := listener.Accept()
		if err != nil {
			log.Print("Connection refused!")
			continue
		}
		defer func() {
			log.Print("Closing connection")
			if connection != nil {
				connection.Close()
			}
		}()
		go handler(connection, command, quit, listener)
	}
}

func gfx(command chan string) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Example #1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	primarySurface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	primarySurface.FillRect(nil, sdl.MapRGB(primarySurface.Format, 192, 255, 192))
	window.UpdateSurface()
	<-command
}

func main() {
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Can't open the port!", err)
	}
	defer l.Close()

	command := make(chan string)

	go gfx(command)
	server(l, command)
}
