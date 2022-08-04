package api

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dennisdebest/joke-fetcher/helper"
	"github.com/emvi/null"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

//go:embed dataset/api.json
var dataset []byte

var ListApis Apis

type Apis struct {
	Apis []Api `json:"apis"`
}

type Api struct {
	Name            string                 `json:"name"`
	Title           string                 `json:"title"`
	Url             string                 `json:"url"`
	Field           null.String            `json:"field"`
	QueryParameters map[string]interface{} `json:"queryParams"`
}

func FetchJoke() string {
	arguments := helper.GetArguments()
	return CallApi(arguments, dataset)
}

func GetApis() Apis {
	var apis Apis
	err := json.Unmarshal(dataset, &apis)
	if err != nil {
		fmt.Println(err)
	}
	return apis
}

func GetApiNames() []string {
	var names []string
	for _, api := range ListApis.Apis {
		names = append(names, api.Name)
	}
	return names
}

func CallApiByName(name string) string {
	apiPtr, err := getApi(name)
	if err != nil {
		fmt.Print(err)
	}
	api := *apiPtr

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
		log.Fatal(response)
		joke = response[field.(string)].(string)
	}

	return joke
}

func CallApi(arguments helper.Arguments, dataset []byte) string {

	name := arguments.Name
	verbose := arguments.Verbose

	apiPtr, err := getApi(name)
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

	client := http.Client{}
	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	req.Header = http.Header{
		"Accept": {"application/json"},
	}

	response, err := client.Do(req)

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

	if verbose {
		fmt.Printf("%v", joke)
	}

	return joke
}

func handleQueryParameters(queryParameters map[string]interface{}, apiUrl *url.URL) {
	params := url.Values{}
	for key, value := range queryParameters {
		params.Add(key, value.(string))
	}
	apiUrl.RawQuery = params.Encode()
}

func getApi(name string) (*Api, error) {
	apis := GetApis().Apis
	names := GetApiNames()

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
