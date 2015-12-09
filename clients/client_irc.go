package clients

import (
	"fmt"
	"log"
	"regexp"

	"github.com/nickvanw/ircx"
	"github.com/snichme/gobot/types"
	"github.com/sorcix/irc"
)

var (
	name     = "LanderblomGoBot"
	server   = "rajaniemi.freenode.net:6667" //"sinisalo.freenode.net:6667"
	channels = "#landerblom"
)

func getGroupForUser(robot types.Robot, username string) string {
	users := robot.Brain().Get("users").([]map[string]string)
	for _, user := range users {
		if user["username"] == username {
			return user["group"]
		}
	}
	return "guest"
}

func send(robot types.Robot, s ircx.Sender, message, target, user string) {
	q := types.Query{
		Statement: message,
		Context: types.QueryContext{
			Username: user,
			Group:    getGroupForUser(robot, user),
		},
	}

	if found, c := robot.Query(q); found {
		for answer := range c {
			fmt.Println(answer)
			s.Send(&irc.Message{
				Command: irc.PRIVMSG,
				Params:  []string{target, ":" + string(answer)},
			})
		}
	}
}

func (client IRCClient) onMessage(s ircx.Sender, m *irc.Message) {

	if m.Params[0] == name { // Private message
		send(client.robot, s, m.Trailing, m.Prefix.Name, m.Prefix.Name)
		return
	}

	reStr := fmt.Sprintf("^@%s ([@\\w\\s]+)", name)
	re := regexp.MustCompile(reStr)
	matches := re.FindStringSubmatch(m.Trailing)
	if matches != nil { // Channel message where bot is mentioned first in message
		send(client.robot, s, matches[1], channels, m.Prefix.Name)
		return
	}
}

func (client IRCClient) onConnected(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channels},
	})
}

//IRCClient The robot joins IRC
type IRCClient struct {
	robot types.Robot
}

// Start start the client
func (client IRCClient) Start() {

	bot := ircx.Classic(server, name)
	if err := bot.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}
	bot.HandleFunc(irc.RPL_WELCOME, client.onConnected)
	bot.HandleFunc(irc.PRIVMSG, client.onMessage)

	fmt.Fprintf(client.robot, "IRCClient: Joined %s channel %s as %s", server, "#landerblom", name)
	bot.HandleLoop()
}

// NewIRCClient Get a new IRC Client
func NewIRCClient(robot types.Robot) *IRCClient {
	return &IRCClient{
		robot: robot,
	}
}
