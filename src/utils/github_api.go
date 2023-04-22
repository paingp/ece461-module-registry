package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Function to get the GitHub URL if input is a NPM URL
func GetGithubUrl(url string) string {
	before, after, found := strings.Cut(url, "www")
	//Finding endpoints and checking for their existence
	if found {
		npmEndpoint := before + "registry" + after
		npmEndpoint = strings.Replace(npmEndpoint, "com", "org", 1)
		npmEndpoint = strings.Replace(npmEndpoint, "package/", "", 1)

		resp, err := http.Get(npmEndpoint)

		if err != nil {
			log.Fatalf("Failed GET request to NPM API at: %s\n", npmEndpoint)
			return ""
		}

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)

			if err != nil {
				log.Fatalf("Failed to read response from NPM API at %s\n", npmEndpoint)
				return ""
			}

			bodyString := string(bodyBytes)
			resBytes := []byte(bodyString)
			var npmRes map[string]interface{}
			_ = json.Unmarshal(resBytes, &npmRes)

			//Checking for existence of GitHub url
			if npmRes["bugs"] == nil {
				return ""
			}

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

func CloneRepo(repo string, directory string) string {

	//GITHUB_TOKEN := os.Getenv("GITHUB_TOKEN")

	lastIdx := strings.LastIndex(repo, "/")
	dir := path.Join(directory, repo[lastIdx+1:])

	err := os.MkdirAll(dir, 0777)

	if err != nil {
		log.Fatalf("Failed to create directory %s to store repo %s\n%s", directory, repo, err)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:          repo + ".git",
		SingleBranch: false,
		Depth:        1,
	})

	if err != nil {
		//metrics.Functions = append(metrics.Functions, "Can't clone "+repo+".git")
		log.Fatalf("Failed to clone repo: %s\n%s", repo, err)
		return ""
	}
	return dir
}

func GetGithubEndpoint(url string) string {
	index := strings.Index(url, "github")
	url = "https://api." + strings.Replace(url[index:], "/", "/repos/", 1)
	return url
}

func GetDataFromGithub(client *http.Client, url string) (map[string]interface{}, error) {
	resp, err := client.Get(url)

	if (err != nil) || (resp.StatusCode != http.StatusOK) {
		log.Fatalf("Failed GET request to GITHUB API at: %s\n", url)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read response body from GITHUB API at %s\n", url)
		return nil, err
	}

	var jsonData map[string]interface{}
	_ = json.Unmarshal(data, &jsonData)

	return jsonData, err
}

func GetPRs(client *http.Client, url string) []map[string]interface{} {
	resp, error := client.Get(url)

	if error != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to GET from GitHub API to compute GoodEngineeringProcess Score\n")
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	bodyString := string(bodyBytes)
	resBytes := []byte(bodyString)
	var npmRes []map[string]interface{}
	_ = json.Unmarshal(resBytes, &npmRes)

	return npmRes
}
