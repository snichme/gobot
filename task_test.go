package main

import (
	"strings"
	"time"
)

type TestTask struct {
}

func (tt TestTask) Name() string {
	return "Test"
}
func (tt TestTask) HelpText() string {
	return "Test test testing..."
}
func (tt TestTask) CanHandle(query Query) bool {
	return strings.Contains(query.Statement, "mange")
}

func (tt TestTask) DoHandle(query Query) <-chan Answer {
	c1 := make(chan Answer, 1)
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "Response from TestTask with delay"
		close(c1)
	}()
	return c1 //query.Client.Respond(<-c1)
}
