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
	joke := api.FetchJoke()
	fmt.Println(joke)
}

func init() {
	rand.Seed(time.Now().Unix())
	helper.DefineArguments()
}
