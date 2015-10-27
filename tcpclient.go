package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type TCPConn struct {
	conn  net.Conn
	robot Robot
}

func newTCPConn(robot Robot, conn net.Conn) {
	c := &TCPConn{
		conn:  conn,
		robot: robot,
	}
	c.Read()
	conn.Close()

}

func (tp TCPConn) Read() {
	for {
		buf := make([]byte, 1024)

		tp.conn.Write([]byte(fmt.Sprintf("%s > ", blue("User"))))
		bytes, err := tp.conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		message := strings.Trim(string(buf[:bytes]), " \r\n")
		if message == "quit" {
			break
		}
		tp.robot.Write(Query{
			Statement: message,
			Client:    tp,
		})
	}
}
func (tp TCPConn) Recieve(in <-chan Answer) {
	for answer := range in {
		tp.conn.Write([]byte(fmt.Sprintf("%s > %s\n", blue(tp.robot.Name), grey(string(answer)))))
	}
}

type TCPClient struct {
	robot Robot
}

func (cc TCPClient) Start() {
	uri := "0.0.0.0:" + cc.robot.config["TCP_PORT"]
	l, err := net.Listen("tcp", uri)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + uri)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go newTCPConn(cc.robot, conn)
	}

}

func NewTCPClient(robot Robot) *TCPClient {
	return &TCPClient{
		robot: robot,
	}
}
