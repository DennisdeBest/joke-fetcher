package api

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emvi/null"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"joker/helper"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

//go:embed dataset/api.json
var dataset []byte

type Apis struct {
	Apis []Api `json:"apis"`
}

type Api struct {
	Name            string                 `json:"name"`
	Url             string                 `json:"url"`
	Field           null.String            `json:"field"`
	QueryParameters map[string]interface{} `json:"queryParams"`
}

func FetchJoke() {
	arguments := helper.GetArguments()
	CallApi(arguments, dataset)
}

func GetApis(dataset []byte) Apis {
	var apis Apis
	err := json.Unmarshal(dataset, &apis)
	if err != nil {
		fmt.Println(err)
	}
	return apis
}

func GetApiNames(apis []Api) []string {
	var names []string
	for _, api := range apis {
		names = append(names, api.Name)
	}
	return names
}

func CallApi(arguments helper.Arguments, dataset []byte) {

	name := arguments.Name
	verbose := arguments.Verbose

	apiPtr, err := getApi(name, dataset)
	if err != nil {
		fmt.Print(err)
	}
	api := *apiPtr

	if verbose {
		fmt.Printf("Api name : %v \n", api.Name)
	}

	apiUrl, err := url.Parse(api.Url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	handleQueryParameters(api.QueryParameters, apiUrl)

	response, err := http.Get(apiUrl.String())

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	field, _ := api.Field.Value()

	var joke = string(responseData)

	if field != nil {
		var response map[string]interface{}
		json.Unmarshal(responseData, &response)
		joke = response[field.(string)].(string)
	}

	Joke = joke

	if verbose {
		fmt.Printf("%v", joke)
	}
}

var Joke string

func GetJoke() string {
	return Joke
}

func handleQueryParameters(queryParameters map[string]interface{}, apiUrl *url.URL) {
	params := url.Values{}
	for key, value := range queryParameters {
		params.Add(key, value.(string))
	}
	apiUrl.RawQuery = params.Encode()
}

func getApi(name string, dataset []byte) (*Api, error) {
	apis := GetApis(dataset).Apis
	names := GetApiNames(apis)

	var api Api
	if name != "" {
		if !slices.Contains(names, name) {
			return nil, errors.New(fmt.Sprintf("%v is not an available API, available ones are : %v\n", name, names))
		} else {
			apiPointer, _ := getApiByName(apis, name)
			api = *apiPointer
		}
	} else {
		api = getRandomApi(apis)
	}

	return &api, nil
}

func getRandomApi(apis []Api) Api {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(apis)
	return apis[n]
}

func getApiByName(apis []Api, name string) (*Api, error) {
	for _, api := range apis {
		if api.Name == name {
			return &api, nil
		}
	}

	return nil, errors.New("no api was found with that name")
}
