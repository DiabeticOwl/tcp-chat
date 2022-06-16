// TCP-Chat
// Small chat using TCP and Go. Interact with other users by stating your
// name first.
package main

import (
	"bufio"
	"fmt"
	"net"
)

// user is a struct that contains the principal information of the User
// of this application.
type user struct {
	UserName string
	Messages chan string
	Conn     *net.Conn
}

// conns is an map of connections that will store every connection made.
// The map value is always true.
var conns = make(map[*net.Conn]bool)

func main() {
	// A listener is opened in the port "8080" with the "tcp" network.
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer li.Close()

	// Infinitely accepts for new connections and sends a goroutine for
	// each one that handles it.
	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}

		go handle(conn)
	}
}

// handle takes the accepted connection and asks for an identification
// from the user. After a scanner brought by the "bufio" package will
// read each line inputted by the user and print them to the server and
// any other connection in the "conns" array.
func handle(conn net.Conn) {
	defer conn.Close()

	// i is a counter used in the scanning of lines for determining
	// in which line of all the messages that the user has sent the
	// scanner is located.
	var u user

	fmt.Fprintf(conn, "Please enter your name: ")

	scanner := bufio.NewScanner(conn)
	scanner.Scan()

	u = user{
		UserName: scanner.Text(),
		Messages: make(chan string),
		Conn:     &conn,
	}

	conns[&conn] = true

	fmt.Println(fmt.Sprintf("User %v has logged in.", u.UserName))

	for scanner.Scan() {
		ln := scanner.Text()
		go func() { u.Messages <- ln }()
		sendingMessage(u)
	}

	delete(conns, &conn)
	fmt.Println(fmt.Sprintf("User %v has logged out.", u.UserName))
}

// sendingMessage will take a message that was set on the "Messages" channel
// for the user and print it to the server and the other connections.
func sendingMessage(u user) {
	var m string

	m = fmt.Sprintf("User %v says: %v", u.UserName, <-u.Messages)

	for c := range conns {
		if c != u.Conn {
			// *() is to dereference the connection for the "fmt" package.
			fmt.Fprintln(*(c), m)
		}
	}

	fmt.Println(m)
}
