package ratom

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Module struct {
	data     packageData
	metadata packageData
}

type packageData struct {
	content   string
	url       string
	jsprogram string
}

type packageMetadata struct {
	name    string
	version string
	id      string
	readme  string
}

func GetToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

// Function to get the GitHub URL from the npmurl input
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

func Clone(repo string) string {

	GITHUB_TOKEN := GetToken()
	log.Printf(GITHUB_TOKEN)

	lastIdx := strings.LastIndex(repo, "/")
	dir := "temp/" + repo[lastIdx+1:]

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

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
