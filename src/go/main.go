package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	// "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	// "fmt"

	// "math"

	"net/http"
	"strings"

	// "utils"

	"github.com/hugoday/ECE461ProjectCLI/src/go/ratom"
)

type NoLog int

func (NoLog) Write([]byte) (int, error) {
	return 0, nil
}


var files []string

func main() {

	doLogging := true
	logFileName := os.Getenv("LOG_FILE")
	logLevel := os.Getenv("LOG_LEVEL")
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil || (logLevel != "1" && logLevel != "2") {
		doLogging = false
	}

	if doLogging {
		if logLevel == "2" {
			ratom.DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		} else {
			ratom.DebugLogger = log.New(new(NoLog), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		}
		ratom.InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		ratom.DebugLogger = log.New(new(NoLog), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		ratom.InfoLogger = log.New(new(NoLog), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	// Makes sure repository folder is clear
	ratom.ClearRepoFolder()

	var inputType int
	fmt.Print("Input type:")
	fmt.Scanln(&inputType)

	// Opens URL file and creates a scanner
	// file, _ := os.Open(os.Args[1])
	// scanner := bufio.NewScanner(file)

	// Create head and temporary repo nodes
	var head *ratom.Repo
	var hold *ratom.Repo
	head = &ratom.Repo{URL: "HEAD"}

	ratom.InfoLogger.Println("Beginning URL file read")

	if inputType == 1 {
		// for each url in the file

		file, _ := os.Open(os.Args[1])
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			//Create new repositories with current URL scanned
			hold = ratom.NewRepo(scanner.Text())
			ratom.InfoLogger.Println("New repo created successfully")
			// Adds repository to Linked List in sorted order by net score
			head = ratom.AddRepo(head, head.Next, hold)
		}
	}
	if inputType == 2 {
		path, _ := os.Getwd()
		cmd := exec.Command("mkdir", path+"/zipTemp")

		// The `Output` method executes the command and
		// collects the output, returning its value
		_, err := cmd.Output()
		if err != nil {
			// if there was any error, print it here
			fmt.Println("could not run command: ", err)
		}

		// reader, err := zip.OpenReader("/Users/adityasrikanth/Desktop/check.zip")
		ratom.Unzip(os.Args[1], path+"/zipTemp")
		// filePath_package := ratom.Walkthrough("/Users/adityasrikanth/Desktop/zipFolder")
		filepath.Walk(path+"/zipTemp", VisitFile)

		filePath_package := files[0]

		jsonFile, err := os.Open(filePath_package)

		if err != nil {
			fmt.Println(err)
		}

		type Data struct {
			Homepage string
			Name     string
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var module_home Data

		err1 := json.Unmarshal(byteValue, &module_home)

		if err1 != nil {
			fmt.Println(err1)
		}

		hold = ratom.NewRepo(module_home.Homepage)
		head = ratom.AddRepo(head, head.Next, hold)

		fmt.Print(head.Next)

		os.RemoveAll(path + "/zipTemp")

		fmt.Print("Name: ", module_home.Name)
		fmt.Print("\nHomepage: ", module_home.Homepage)

		ratom.SetMetadata("tomr-bucket", module_home.Name, head.Next)
	}

	// Prints each repository in NDJSON format to stdout (sorted highest to low based off net score)
	// ratom.PrintRepo(head.Next)

}

func VisitFile(path string, info os.FileInfo, err error) error {

	if err != nil {

		fmt.Println(err)
		return nil
	}

	if info.IsDir() || filepath.Ext(path) != ".json" {

		return nil
	}

	reg, err2 := regexp.Compile("^package")

	if err2 != nil {

		return err2
	}

	if reg.MatchString(info.Name()) {

		files = append(files, path)
	}

	return nil
}

// func main() {
//ratom.ClearRepoFolder()
// 	fmt.Println("Enter input type: ")
// 	var inputType int

// 	var head *ratom.Repo
// 	head = &ratom.Repo{URL: "HEAD"}

// 	// Taking input from user
// 	fmt.Scanln(&inputType)

// 	if inputType == 1{
// 		fmt.Println("Enter input type: ")
// 		var url string
// 		fmt.Scanln(&url)

// 		head = handle_url(url, head)
// 		ratom.PrintRepo(head.Next)
// 	}
// }

// func handle_url(url string, head *ratom.Repo) *ratom.Repo {
// 	var hold *ratom.Repo

// 	url = getGithubUrl(url)

// 	hold = ratom.NewRepo(url)
// 	// head = ratom.AddRepo(head, head.Next, hold)
// 	fmt.Print(hold)

// 	return head
// }

// Function to get the GitHub URL from the npmurl input
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

			if npmEndpoint == "" {
				return ""
			}

			url = strings.Replace(npmEndpoint, "/issues", "", 1)
		}
	}
	return url
}
