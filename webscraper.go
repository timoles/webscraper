package main

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"io"
	"sync"
	"encoding/json"
)

type Config struct{
	Threads int
	ResponseFilePath string
}

var configPath string = "./config.conf"
var config Config



func main() {
	//https://www.zerodayinitiative.com/advisories/ZDI-17-001/
	// Check args
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]
	//url = "127.0.0.1/%d"
	f, err := ioutil.ReadFile(configPath)
	// TODO make function for error handling
	if err != nil{
		fmt.Println("Config file not found exiting...")
		os.Exit(1)
	}
	configJson := string(f)

	json.Unmarshal([]byte(configJson), &config)

	// Create directory into which we will save the responses
	if _, err := os.Stat(config.ResponseFilePath); os.IsNotExist(err) {
		os.Mkdir(config.ResponseFilePath, os.ModePerm)
	}

	var wg sync.WaitGroup

	queryResponseChannel := make(chan bool, 3)
	//for i := 1010; i < 1015 && missingCount < 2; i++ {
	for i := 0; i < 10; i++ {
		//missingCount := 0
		wg.Add(1)
		go querySite(queryResponseChannel, &wg, url, i)
	}
	//<-queryResponseChannel
	wg.Wait()
	close(queryResponseChannel)

}

func querySite(queryResponseChannel chan<- bool, wg *sync.WaitGroup, queryUrl string, i int) {
	defer wg.Done()
	urlEdited := fmt.Sprintf(queryUrl, i)
	fmt.Println("Sende Request für: " + urlEdited)
	response, err := http.Get(urlEdited)
	responseDataString := ""
	if response == nil {
		fmt.Println("Response nil, error")
		return
	}
	if response.Status != "200 OK" {
		fmt.Println("Unknown Status: " + response.Status + " " + urlEdited)
	} else {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		responseDataString = string(responseData)

		//if strings.Contains(responseDataString, "advisories-details") && !(responseDataStringOld == responseDataString) { // TODO make real xml parser and only write important tags
		if strings.Contains(responseDataString, "advisories-details") { // TODO make real xml parser and only write important tags
			fmt.Println("Writing: " + urlEdited)

		} else {
			fmt.Println("Duplicate/not expected: " + urlEdited)

		}
		//responseDataStringOld = responseDataString

		response.Body.Close()
	}

	if err != nil {
		log.Fatal(err)
	}
	go writeToDisk(responseDataString, urlEdited)

	//queryResponseChannel <- true
}

func writeToDiskGo(writeChannel chan<- bool, responseBodyString, queryUrlEdited string) {
	fmt.Println("Schreibe Datei für; " + queryUrlEdited)
	index := strings.LastIndex(queryUrlEdited[:len(queryUrlEdited)-1], "/")
	responseIndex := queryUrlEdited[index:]
	responseIndex = strings.Replace(responseIndex, "/", "", -1) //-1 unlimited replacements

	f, err := os.Create(config.ResponseFilePath + responseIndex + ".html")
	defer f.Close()
	io.WriteString(f, responseBodyString)
	if err != nil {
		log.Fatal(err)
	}
	writeChannel <- true
}
func writeToDisk(responseBodyString, queryUrlEdited string) {
	fmt.Println("Schreibe Datei für; " + queryUrlEdited)
	index := strings.LastIndex(queryUrlEdited[:len(queryUrlEdited)-1], "/")
	responseIndex := queryUrlEdited[index:]
	responseIndex = strings.Replace(responseIndex, "/", "", -1) //-1 unlimited replacements

	f, err := os.Create(config.ResponseFilePath + responseIndex + ".html")
	defer f.Close()
	io.WriteString(f, responseBodyString)
	if err != nil {
		log.Fatal(err)
	}
}
