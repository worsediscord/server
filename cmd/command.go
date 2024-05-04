package main

type Command interface {
	Name() string
	Description() string
	Parse(args []string) error
	Run() error
}
