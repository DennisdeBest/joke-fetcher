package main

import (
	_ "embed"
	"fmt"
	"github.com/dennisdebest/joke-fetcher/api"
	"github.com/dennisdebest/joke-fetcher/helper"
	"math/rand"
	"time"
)

func main() {
	joke, err := api.FetchJoke()
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(joke)
}

func init() {
	rand.Seed(time.Now().Unix())
	helper.DefineArguments()
}
