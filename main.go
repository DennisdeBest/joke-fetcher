package main

import (
	_ "embed"
	"fmt"
	"github.com/dennisdebest/joke-fetcher/api"
	"github.com/dennisdebest/joke-fetcher/helper"
)

func main() {
	api.FetchJoke()
	fmt.Println(api.Joke)
}

func init() {
	helper.DefineArguments()
}
