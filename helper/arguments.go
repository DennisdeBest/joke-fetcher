package helper

import (
	"flag"
)

type Arguments struct {
	Name    string
	Verbose bool
}

var name string
var verbose bool

func DefineArguments() {
	flag.StringVar(&name, "name", "", "Name of the api to call")
	flag.BoolVar(&verbose, "verbose", false, "Add more output data")
}

func GetArguments() Arguments {
	flag.Parse()
	return Arguments{Name: name, Verbose: verbose}
}
