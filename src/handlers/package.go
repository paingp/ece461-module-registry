package handlers

import (
	"fmt"
	"log"
	"path"
	"tomr/models"
	"tomr/src/metrics"
	"tomr/src/utils"
)

const pkgDirPath = "src/metrics/temp" // temp directory to store packages

func CreatePackage(content string, url string, jsprogram string) {
	packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
	//utils.PrintPackageData(packageData)
	// Return Error 400 if both Content and URL are set
	if (packageData.Content != "") && (packageData.URL != "") {
		fmt.Printf("Error 400: Content and URL cannot be both set")
	} else if packageData.Content != "" { // Only Content is set
		// Decode base64 string into zip
		pkgDir := path.Join(pkgDirPath, "package.zip")
		utils.Base64ToZip(packageData.Content, pkgDir)
		var readMe []byte
		metadata, err := utils.GetMetadataFromZip(pkgDir, &readMe)
		if err != nil {
			log.Fatalf("Failed to get metadata from zip file\n")
		}
		utils.PrintMetadata(metadata)
		metrics.RatePackage(metadata.RepoURL, pkgDir, metadata.License, &readMe)
	} else {
		gitUrl := utils.GetGithubUrl(url)
		pkgDir := utils.CloneRepo(gitUrl, pkgDirPath)
		metrics.RatePackage(gitUrl, pkgDir, "", nil)
	}
	// Upload package and store data in system
}

/*
func GetPackageMetadata1(directory string, isZip bool) {
	pkgJsonPath := ""
	var rc io.ReadCloser
	//metadata := models.PackageMetadata{}
	if isZip {
		pkgJsonPath = utils.GetMetadataFromZip(directory)
	} else {
		pkgJsonPath = path.Join(pkgJsonPath, "/package.json")
	}
	rc, err := os.Open(pkgJsonPath)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(rc)
	type metadata struct {
		Name    string `json:"Name"`
		Version string `json:"Version"`
		License string `json:"License"`
	}
	var m metadata
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	rc.Close()
	fmt.Println(m.Name)
	fmt.Println(m.Version)
	fmt.Println(m.License)

		data, err := os.ReadFile(pkgJsonPath)
		if err != nil {
			log.Fatal(err)
		}

		//packageMetadata := models.PackageMetadata{}
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(data), &jsonMap)

		metadata := models.PackageMetadata{Name: jsonMap["name"].(string), Version: jsonMap["version"].(string), ID: ""}
		metadata.ID = metadata.Name + "_" + metadata.Version

	//utils.PrintMetadata(metadata)
}
*/
