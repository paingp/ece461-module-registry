package utils

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"tomr/models"

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

	GITHUB_TOKEN := os.Getenv("GITHUB_TOKEN")
	log.Printf(GITHUB_TOKEN)

	lastIdx := strings.LastIndex(repo, "/")
	dir := directory + repo[lastIdx+1:]

	err := os.MkdirAll(dir, 0777)

	if err != nil {
		log.Fatal(err)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:          repo + ".git",
		SingleBranch: false,
		Depth:        1,
	})

	if err != nil {
		//metrics.Functions = append(metrics.Functions, "Can't clone "+repo+".git")
		log.Fatal(err)
		return "err"
	}
	return dir
}

func GetMetadataFromZip(zipFile string, metadata *models.PackageMetadata) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "package.json") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			dec := json.NewDecoder(rc)
			for {
				if err := dec.Decode(metadata); err == io.EOF {
					break
				} else if err != nil {
					return err
				}
			}
			rc.Close()
		}
	}
	return nil
}

func DecodeBase64(b64string string) string {
	data, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	err = os.WriteFile("temp/pkg", data, 0777)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return ""
}

func zip2base64(zipFile string) string {
	//zipFile := "lodash.zip"
	bytes, err := os.ReadFile(zipFile)

	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
