package utils

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"tomr/models"
)

func ZipToBase64(zipFile string) (string, error) {
	bytes, err := os.ReadFile(zipFile)

	if err != nil {
		log.Fatalf("Failed to read ZIP file: %s\n%s", zipFile, err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), err
}

func Base64ToZip(b64string string, zipDirectory string) error {
	data, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		log.Fatalf("Failed to decode base64 string into Zip file: %s\n", b64string)
		return err
	}

	err = os.WriteFile(zipDirectory, data, 0777)
	if err != nil {
		log.Fatalf("Failed to write base64 string into the following destination: %s\n", zipDirectory)
		return err
	}
	return err
}

func GetPackageMetadata(directory string) models.PackageMetadata {
	metadata := models.PackageMetadata{}

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

	metadata.ID = metadata.Name + "_" + metadata.Version
	//PrintMetadata(metadata)
	return metadata
}

func GetMetadataFromZip(zipFile string, readme *[]byte) (models.PackageMetadata, error) {
	metadata := models.PackageMetadata{}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return metadata, err
	}
	defer r.Close()
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "package.json") {
			rc, err := f.Open()
			if err != nil {
				return metadata, err
			}
			dec := json.NewDecoder(rc)
			for {
				if err := dec.Decode(&metadata); err == io.EOF {
					break
				} else if err != nil {
					return metadata, err
				}
			}
			rc.Close()
		} else if strings.Count(f.Name, "/") == 1 {
			matched, _ := regexp.MatchString(`(?i)readme`, f.Name)
			if matched {
				fmt.Printf("Matched: %s\n", f.Name)
				*readme = GetReadMeFromZip(f)
			}
		}
	}
	//fmt.Printf(string(readme))
	return metadata, err
}

func GetReadMeFromZip(readme *zip.File) []byte {
	rc, err := readme.Open()
	if err != nil {
		log.Fatalf("Failed to open ReadMe file: %s", readme.Name)
		log.Fatal(err)
		return nil
	}

	bytes, err := io.ReadAll(rc)
	if err != nil {
		log.Fatalf("Failed to read %s", readme.Name)
		log.Fatal(err)
		return nil
	}
	return bytes
}

/*
func walk(path string, d fs.DirEntry, err error) error {
	maxDepth := 1
	if (err != nil) {
		return err
	}
	if d.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxDepth {
		return fs.SkipDir
	} else {
		// Checking paths
		matched, _ := regexp.MatchString(`(?i)readme`, path)
		if (matched) {
			// Checking matched path
			check, _ := regexp.MatchString("(?i)guid", path)
			if !check {
				// Finding readme
				readme = path
			}
		}
	}
	return nil
}
*/
