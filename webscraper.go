package main

import (
"fmt"

"log"
"net/http"
"os"
	"strings"
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
	attachment := 1
	for i := 999; i < 1001; i++ {

		attachment += i
	}
	url = "https://www.zerodayinitiative.com/advisories/ZDI-17-%03v/"
	//number := strconv.Itoa(0)
	//number += strconv.Itoa(0)
	//number += strconv.Itoa(1)

	url = fmt.Sprintf(url,001)
	//url =  fmt.Sprintf(url, number)
	fmt.Println(url)
	return


	//fmt.Println(url[len(url)-5 :])
	//fmt.Println(strings.Count(url, "/"))
	index := strings.LastIndex(url[:len(url)-1], "/")
	responseIndex := url[index:]
	responseIndex = strings.Replace(responseIndex, "/","",-1) //-1 unlimited replacements
	fmt.Println(responseIndex)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		f, err := os.Create(path + responseIndex + ".html")
		defer f.Close()
		defer response.Body.Close()
		response.Write(f)
		if err != nil {
			log.Fatal(err)
		}
	}

}






