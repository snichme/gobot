package clients

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snichme/gobot/types"
)

type CliClient struct {
	robot types.Robot
}

func (cc CliClient) Start() {
	reader := bufio.NewReader(os.Stdin)
	queryContext := types.QueryContext{
		Username: "CLI User",
		Group:    "admin", // Cli user are admin becuase it runs on the same machine as the bot
	}
	queryRobot := func(message string) {
		q := types.Query{
			Statement: message,
			Context:   queryContext,
		}
		if found, c := cc.robot.Query(q); found {
			for answer := range c {
				fmt.Fprint(cc.robot, answer) //, a ...interface{})fmt.Printf("%s > %s\n", blue(cc.robot.Name()), grey(string(answer)))
			}
		}
	}
	for {
		fmt.Printf("%s > ", blue("Admin"))
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, " \n")
		if text == "quit" {
			break
		}
		queryRobot(text)
	}

}

func NewCliClient(robot types.Robot) *CliClient {
	return &CliClient{
		robot: robot,
	}
}
