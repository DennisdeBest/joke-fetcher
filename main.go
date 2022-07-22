package main

import (
	_ "embed"
	"fmt"
	"joker/api"
	"joker/helper"
)

func main() {
	api.FetchJoke()
	fmt.Println(api.Joke)
}

func init() {
	helper.DefineArguments()
}
