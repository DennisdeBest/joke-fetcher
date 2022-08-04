package main

import (
	_ "embed"
	"fmt"
	"github.com/dennisdebest/joke-fetcher/api"
	"github.com/dennisdebest/joke-fetcher/helper"
)

func main() {
	joke := api.FetchJoke()
	fmt.Println(joke)
}

func init() {
	helper.DefineArguments()
}
