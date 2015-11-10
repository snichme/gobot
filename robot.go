package main

import (
	"fmt"
	"os"
	"strings"
)

type RobotConfig struct {
	Settings    map[string]string   `json:"settings"`
	Permissions map[string][]string `json:"permissions"`
}

type RobotBrain interface {
	Get(key string) interface{}
	Set(key string, value interface{}) bool
}

type Robot struct {
	tasks       []Task
	settings    map[string]string
	permissions map[string][]string
	Brain       RobotBrain
}

func (robot Robot) Name() string {
	return robot.settings["name"]
}

func (r Robot) HasAccess(taskName string, context QueryContext) bool {
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
	return fmt.Fprintf(os.Stdout, "%s", p)
}

func (l Robot) Query(s Query) (bool, <-chan Answer) {
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
func (r Robot) CanHandle(query Query) bool {
	return strings.Contains(query.Statement, "help me")
}
func (r Robot) DoHandle(query Query) <-chan Answer {
	c1 := make(chan Answer)
	go func() {
		for _, task := range r.tasks {
			c1 <- Answer(fmt.Sprintf("[%s] %s", task.Name(), task.HelpText()))
		}
		close(c1)
	}()
	return c1
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

func NewRobot(config RobotConfig, tasks []Task) Robot {
	brain := &InMemoryBrain{
		storage: make(map[string]interface{}),
	}

	brain.Set("users", []map[string]string{
		{"username": "mange", "password": "secret", "group": "admin"},
		{"username": "apa", "password": "banan", "group": "member"},
	})

	return Robot{
		tasks:       tasks,
		settings:    config.Settings,
		permissions: config.Permissions,
		Brain:       brain,
	}
}
