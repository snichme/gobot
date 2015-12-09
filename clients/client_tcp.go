package clients

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/snichme/gobot/types"
)

//TCPClient Connect to the robot using TCP
type TCPClient struct {
	robot types.Robot
}

// Start start the client
func (cc TCPClient) Start() {
	uri := "127.0.0.1:" + cc.robot.Setting("tcp_port")
	l, err := net.Listen("tcp", uri)
	if err != nil {
		fmt.Fprintf(cc.robot, "Error listening: %s", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Fprintf(cc.robot, "TcpClient: Listening on %s", uri)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		for {
			buf := make([]byte, 1024)
			conn.Write([]byte(fmt.Sprintf("%s > ", blue("User"))))
			bytes, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}
			message := strings.Trim(string(buf[:bytes]), " \r\n")
			if message == "quit" {
				break
			}
			q := types.Query{
				Statement: message,
				Context: types.QueryContext{
					Username: "TcpUser-1",
					Group:    "guest", // All Tcp users are Guest since no auth
				},
			}
			if found, c := cc.robot.Query(q); found {
				i := 0
				for answer := range c {
					var msg string
					if i == 0 {
						msg = fmt.Sprintf("%s > %s\n", blue(cc.robot.Name()), grey(string(answer)))
					} else {
						format := fmt.Sprintf("%%%ds %%s\n", len(cc.robot.Name())+2)
						msg = fmt.Sprintf(format, " ", grey(string(answer)))
					}
					conn.Write([]byte(msg))
					i++
				}
			}
		}
		conn.Close()
	}
}

// NewTCPClient Get a new TCP Client
func NewTCPClient(robot types.Robot) *TCPClient {
	return &TCPClient{
		robot: robot,
	}
}
