package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetGithubEndpoint(url string) string {
	index := strings.Index(url, "github")
	url = "https://api." + strings.Replace(url[index:], "/", "/repos/", 1)
	return url
}

func GetDataFromGithub(client *http.Client, url string) map[string]interface{} {
	resp, error := client.Get(url)

	if (error != nil) || (resp.StatusCode != http.StatusOK) {
		fmt.Println("HTTP Client get failed")
		return nil
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	var jsonData map[string]interface{}
	_ = json.Unmarshal(data, &jsonData)

	return jsonData
}

func GetPRs(client *http.Client, url string) []map[string]interface{} {
	resp, error := client.Get(url)

	if error != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP Client get failed")
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
