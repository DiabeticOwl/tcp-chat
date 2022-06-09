package main

import (
	"bufio"
	"fmt"
	"net"
)

type user struct {
	UserName string
	Messages chan string
	Conn     *net.Conn
}

var conns []net.Conn

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	var i int = 0
	var u user

	fmt.Fprintf(conn, "Please enter your name: ")

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		ln := scanner.Text()

		if i == 0 {
			u = user{
				UserName: ln,
				Messages: make(chan string),
				Conn:     &conn,
			}

			conns = append(conns, conn)

			fmt.Println(fmt.Sprintf("User %v has logged in.", u.UserName))
		} else {
			go func() { u.Messages <- ln }()
			sendingMessage(u)
		}

		i++
	}

	fmt.Println(fmt.Sprintf("User %v has logged out.", u.UserName))
}

func sendingMessage(u user) {
	var m string

	m = fmt.Sprintf("User %v says: %v", u.UserName, <-u.Messages)

	for _, c := range conns {
		if c != *(u.Conn) {
			fmt.Fprintln(c, m)
		}
	}

	fmt.Println(m)
}
