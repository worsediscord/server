package main

import "reflect"

type callable struct {
	fn   interface{}
	args []interface{}
}

type Cleaner struct {
	callables []callable
}

func (c *Cleaner) RegisterFunc(fn interface{}, args ...interface{}) {
	c.callables = append(c.callables, callable{fn: fn, args: args})
}

func (c *Cleaner) Execute() {
	for _, callable := range c.callables {
		fn := reflect.ValueOf(callable.fn)

		var args []reflect.Value
		for _, arg := range callable.args {
			args = append(args, reflect.ValueOf(arg))
		}

		fn.Call(args)
	}
}
