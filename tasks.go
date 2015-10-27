package main

type Answer string

type Task interface {
	Name() string
	HelpText() string
	CanHandle(query Query) bool
	DoHandle(query Query) <-chan Answer
}
