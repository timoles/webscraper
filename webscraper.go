package main

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"io"
)

func main() {
	//https://www.zerodayinitiative.com/advisories/ZDI-17-001/
	// Check args
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]
	// Create directory into which we will save the responses
	path := "./responses/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	urlEdited := ""
	missingCount := 0
	responseDataStringOld := ""
	for i := 1010; i < 1020 && missingCount < 2; i++ { // TODO multi threaded
		urlEdited = fmt.Sprintf(url, i)

		response, err := http.Get(urlEdited)
		if response.Status != "200" {
			fmt.Println("Unknown Status: " + response.Status + " " + urlEdited)
		} else {

			if err != nil {
				log.Fatal(err)
			} else {
				responseData, err := ioutil.ReadAll(response.Body)
				responseDataString := string(responseData)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(responseDataString, "advisories-details") && !(responseDataStringOld == responseDataString) { // TODO make real xml parser and only write important tags
					fmt.Println("Writing: " + urlEdited)
					missingCount = 0
					index := strings.LastIndex(urlEdited[:len(urlEdited)-1], "/")
					responseIndex := urlEdited[index:]
					responseIndex = strings.Replace(responseIndex, "/", "", -1) //-1 unlimited replacements
					f, err := os.Create(path + responseIndex + ".html")
					io.WriteString(f, responseDataString)
					if err != nil {
						log.Fatal(err)
					}
					f.Close()
				} else {
					fmt.Println("Duplicate/not expected: " + urlEdited)
					missingCount++
				}
				responseDataStringOld = responseDataString
			}
			response.Body.Close()
		}
	}

}
