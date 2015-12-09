package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/snichme/gobot/types"
)

func blue(v string) string {
	return "\033[34m" + v + "\033[0m"
}
func grey(v string) string {
	return "\033[37m" + v + "\033[0m"
}

type Robot struct {
	log         io.Writer
	tasks       []types.Task
	settings    map[string]string
	permissions map[string][]string
	brain       types.RobotBrain
}

func (robot Robot) Name() string {
	return robot.settings["name"]
}

func (r Robot) HasAccess(taskName string, context types.QueryContext) bool {
	contextGroup := context.Group
	groups, ok := r.permissions[taskName]
	if contextGroup == "admin" {
		return true
	}
	if !ok {
		return true
	}
	for _, group := range groups {
		if group == contextGroup {
			return true
		}
	}
	return false
}

func (r Robot) Write(p []byte) (n int, err error) {
	fmt.Fprintf(r.log, "[%v] %s\n", time.Now(), p)
	return fmt.Fprintf(os.Stdout, "%s > %s\n", blue(r.Name()), grey(string(p)))
}

func (l Robot) Query(s types.Query) (bool, <-chan types.Answer) {
	if l.CanHandle(s) {
		return true, l.DoHandle(s)
	}
	for _, task := range l.tasks {
		if l.HasAccess(task.Name(), s.Context) {
			if ok, c := task.Handle(s); ok {
				return true, c
			}
		}
	}
	return false, nil
}

func (r Robot) HelpText() string {
	return "Will help you when in needs"
}
func (r Robot) CanHandle(query types.Query) bool {
	return strings.Contains(query.Statement, "help me")
}
func (r Robot) DoHandle(query types.Query) <-chan types.Answer {
	c1 := make(chan types.Answer)
	go func() {
		for _, task := range r.tasks {
			c1 <- types.Answer(fmt.Sprintf("[%s] %s", task.Name(), task.HelpText()))
		}
		close(c1)
	}()
	return c1
}
func (r Robot) Brain() types.RobotBrain {
	return r.brain
}

func (r Robot) Setting(key string) string {
	return r.settings[key]
}

type InMemoryBrain struct {
	storage map[string]interface{}
}

func (brain InMemoryBrain) Get(key string) interface{} {
	return brain.storage[key]
}
func (brain InMemoryBrain) Set(key string, value interface{}) bool {
	brain.storage[key] = value
	return true
}

func NewRobot(config types.RobotConfig, tasks []types.Task) types.Robot {
	brain := &InMemoryBrain{
		storage: make(map[string]interface{}),
	}

	brain.Set("users", []map[string]string{
		{"username": "mange", "password": "secret", "group": "admin"},
		{"username": "apa", "password": "banan", "group": "member"},
	})

	fo, err := os.Create("robotlog.txt")
	if err != nil {
		log.Fatal(err)
	}
	return Robot{
		log:         fo,
		tasks:       tasks,
		settings:    config.Settings,
		permissions: config.Permissions,
		brain:       brain,
	}
}
