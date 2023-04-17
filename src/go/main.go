package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

		reader, err := zip.OpenReader("/Users/adityasrikanth/Desktop/check.zip")
		if err != nil {
			fmt.Print("errrr")
			return
		}
		defer reader.Close()

		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		for _, f := range reader.File {
			err := unzipFile(f, path)
			if err != nil {
				fmt.Print("2")
				return
			}
		}
	}

	// Prints each repository in NDJSON format to stdout (sorted highest to low based off net score)
	// ratom.PrintRepo(head.Next)
	ratom.SetMetadata("tmr-bucket", "lodash.txt", head.Next)
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
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
