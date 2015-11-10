package main

type Answer string

type Task interface {
	Name() string
	HelpText() string
	Handle(query Query) (bool, <-chan Answer)
}
