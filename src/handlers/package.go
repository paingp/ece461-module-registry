package handlers

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"tomr/models"
	"tomr/src/utils"
)

func GetPackageMetadata(directory string) {
	metadata := models.PackageMetadata{}
	if strings.HasSuffix(directory, ".zip") {
		utils.GetMetadataFromZip(directory, &metadata)
	} else {
		pkgJsonPath := path.Join(directory, "package.json")
		file, err := os.Open(pkgJsonPath)
		if err != nil {
			panic(err)
		}
		dec := json.NewDecoder(file)
		for {
			if err := dec.Decode(&metadata); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}
		file.Close()
	}
	metadata.ID = metadata.Name + "_" + metadata.Version
	utils.PrintMetadata(metadata)
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
