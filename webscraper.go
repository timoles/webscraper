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
	"time"
)
import (
	"github.com/goware/urlx"
)

type Config struct{
	Threads int
	ResponseFilePath string
	Keyword string
}

type Node struct{
	Value string
	Left *Node
	Right *Node
}




var configPath string = "./config.conf"
var config = &Config{1, "hi", "du"}

func main() {
	// TODO checken wie es mit Datenbanklösungen aussieht
	// TODO check if correct tree
	// Todo evtl balanced binary tree, but dont think thats that great with spidering data set
	// TODO way better error handling

	// anchor := &Node{"",nil,nil}
	// fmt.Println(anchor.String())

	//https://www.zerodayinitiative.com/advisories/ZDI-17-%03v/
	// Check args
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]

	url = "https://www.swisscom.ch/de/privatkunden.html" // TODO
	// url = "https://www.zerodayinitiative.com/advisories/ZDI-18-1000/"

	//url = "127.0.0.1/%d"
	f, err := ioutil.ReadFile(configPath)
	if err != nil{
		fmt.Println("Config file not found exiting...")
		log.Fatal(err)
	}
	configJson := string(f)
	err = json.Unmarshal([]byte(configJson), &config)
	if  err != nil{
		fmt.Println("Config file Error exiting...")
		log.Fatal(err)
	}
	fmt.Println(config.Threads)
	// Create directory into which we will save the responses
	if _, err := os.Stat(config.ResponseFilePath); os.IsNotExist(err) {
		os.Mkdir(config.ResponseFilePath, os.ModePerm)
	}

	var wg sync.WaitGroup
	var client = &http.Client{Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		} } // TODO need client for every single query? prob not?
	// TODO check all client options UserAgent etc...
	queryResponseChannel := make(chan bool, 3)
	//for i := 1010; i < 1015 && missingCount < 2; i++ {
	for i := 0; i < 1; i++ {
		//missingCount := 0
		wg.Add(2) // Add 2 one for query and one for writing
		go querySite(queryResponseChannel, &wg, url, i,client)

	}
	//<-queryResponseChannel
	wg.Wait()
	close(queryResponseChannel)

}

func querySite(queryResponseChannel chan<- bool, wg *sync.WaitGroup, queryUrl string, i int, client *http.Client) {
	defer wg.Done()
	// Response
	fmt.Println()
	fmt.Println("Sende Request für: " + queryUrl)

	response, err := client.Get(queryUrl)
	responseDataString := ""
	if err != nil {
		log.Fatal(err)
		return
	}

	if response == nil {
		fmt.Println(" Response nil, error")
		return
	}
	if response.Status != "200 OK" {
		fmt.Println(" Unknown Status: " + response.Status + " " + queryUrl)
	} else {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		responseDataString = string(responseData)
		response.Body.Close()

		queryForUrls(responseDataString,config.Keyword )
	}
	if err != nil {
		log.Fatal(err)
	}
	if responseDataString != ""{
		go writeToDisk(responseDataString, queryUrl, wg) // TODO check if das schon so richtig schnellstmöglich ausführt oder immer zuerst querys und dann schreibe
	}else {
		wg.Done() // Done because writing cant do it
	}
	//queryResponseChannel <- true
}
func queryForUrls(responseData string, keyword string){ // TODO multiple keywords
	//fmt.Println("Query Data with keyword: " + keyword)
	leadingIdentifier := "https://"
	trailingIdentifier := ".html" // TODO shitty identifier
	indexKeyword := strings.Index(responseData, keyword)
	globalDocIndex := 0
	//fmt.Println(indexKeyword)
	indexLeadingIdentifier := -1
	indexTrailingIdentifier := -1
	for ; indexKeyword != -1;{
		//fmt.Println("Found keyword occurence at Index: " ,indexKeyword)
		//fmt.Println(responseData[indexKeyword:])
		indexLeadingIdentifier = strings.LastIndex(responseData[:indexKeyword+len(keyword)], leadingIdentifier)
		indexTrailingIdentifier = strings.Index(responseData[indexKeyword:], trailingIdentifier) + indexKeyword // TODO shitty identifier
		if indexLeadingIdentifier != -1 && indexTrailingIdentifier != -1{
			//fmt.Println(responseData)
			//fmt.Println(indexLeadingIdentifier, " " , indexTrailingIdentifier)
			//fmt.Println("Occurance was valid")
			uriToCheck := responseData[indexLeadingIdentifier:indexTrailingIdentifier]

			if isValidUri(uriToCheck){
				//fmt.Println("Valid Uri found: " + uriToCheck)
				fmt.Println("TODO") // TODO
			}
		}
		globalDocIndex += indexTrailingIdentifier-len(trailingIdentifier)
		indexKeyword = strings.Index(responseData[globalDocIndex:], keyword)
		// TODO href="//tags.tiqcdn.com
		// and prob many more
	}
	fmt.Println("Query done")
}

func isValidUri(toCheck string)bool{
	//fmt.Println("------------------------------")
	//fmt.Println(toCheck)
	//_, err := url.ParseRequestURI(toCheck)
	url , err := urlx.Parse(toCheck)
	fmt.Println(url)
	// fmt.Println(reflect.TypeOf(url))
	// normalized, _ := urlx.Normalize(url)
	// fmt.Println(normalized)
	if err != nil {
		return false
	} else {
		return true
	}
}

func writeToDisk(responseBodyString, queryUrlEdited string, wg *sync.WaitGroup) {
	defer wg.Done()
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
func compareUrl1Smaller(url1 string, url2 string) bool{
	url1 = strings.ToLower(url1) // TODO case sensitive
	url2 = strings.ToLower(url2)
	var length int
	if len(url1) < len(url2){
		length = len(url1)
	}else{
		length = len(url2)
	}
	if length <=0{
		return false
	}
	for i := 0; i < length; i++{
		if url1[i]>url2[i]{
			return false
		}
	}
	return true
}

// Binary Tree
func insert(current *Node, value string) (bool, *Node){
	duplicate := true
	if current == nil{
		return false, &Node{value,nil,nil}
	}else if current.Value == value{
		return true,current
	}
	if compareUrl1Smaller(value, current.Value){
		duplicate, current.Left = insert(current.Left,value)
	}else{
		duplicate, current.Right = insert(current.Right,value)
	}

	return duplicate, current // TODO
}

// To string taken from https://github.com/golang/tour/blob/master/tree/tree.go#L20
func (t *Node) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}
//Binar Tree end