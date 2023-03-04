package main

import (
	"bufio"
	// "fmt"
	"log"
	// "math"
	// "utils"
	"fmt"
	"os"
	"encoding/json"

	//"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
)

type NoLog int

func (NoLog) Write([]byte) (int, error) {
	return 0, nil
}

func main() {

	fmt.Println("Aditya Srikanth")
	fmt.Println("Here now")

	doLogging := true
	logFileName := os.Getenv("LOG_FILE")
	logLevel := os.Getenv("LOG_LEVEL")
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil || (logLevel != "1" && logLevel != "2") {
		doLogging = false
	}

	if doLogging {
		if logLevel == "2" {
			DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		} else {
			DebugLogger = log.New(new(NoLog), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		}
		InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		DebugLogger = log.New(new(NoLog), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		InfoLogger = log.New(new(NoLog), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Makes sure repository folder is clear
	clearRepoFolder()

	// Opens URL file and creates a scanner
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)

	// Create head and temporary repo nodes
	var head *repo
	var hold *repo
	head = &repo{URL: "HEAD"}

	InfoLogger.Println("Beginning URL file read")
	// for each url in the file
	for scanner.Scan() {
		//Create new repositories with current URL scanned
		var tempURL = scanner.Text()
		tempURL = getGithubUrl(tempURL)
		hold = newRepo(tempURL)
		fmt.Println(tempURL)
		fmt.Println("Here now\n\n\n")
		InfoLogger.Println("New repo created successfully")
		// Adds repository to Linked List in sorted order by net score
		head = addRepo(head, head.next, hold)
	}

	// Prints each repository in NDJSON format to stdout (sorted highest to low based off net score)
	printRepo(head.next)
}


func getGithubUrl(url string) string {
	before, after, found := strings.Cut(url, "www")
	//Finding endpoints and checking for their existence
	if found {
		npmEndpoint := before + "registry" + after
		npmEndpoint = strings.Replace(npmEndpoint, "com", "org", 1)
		npmEndpoint = strings.Replace(npmEndpoint, "package/", "", 1)

		resp, err := http.Get(npmEndpoint)

		if err != nil {
			return ""
		}

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)

			if err != nil {
				return ""
			}

			bodyString := string(bodyBytes)
			resBytes := []byte(bodyString)
			var npmRes map[string]interface{}
			_ = json.Unmarshal(resBytes, &npmRes)
			

			bugs := npmRes["bugs"].(map[string]interface{})
			npmEndpoint = bugs["url"].(string)

			if (npmEndpoint == ""){
				return ""
			}

			url = strings.Replace(npmEndpoint, "/issues", "", 1)
		}
	}
	return url
}