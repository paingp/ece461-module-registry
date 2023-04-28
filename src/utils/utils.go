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
	"path/filepath"
	"regexp"
	"strings"
	"tomr/models"
)

func ZipToBase64(zipFile string) ([]byte, error) {
	bytes, err := os.ReadFile(zipFile)

	if err != nil {
		log.Fatalf("Failed to read ZIP file: %s\n%s", zipFile, err)
		return nil, err
	}
	dest := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(dest, bytes)

	return dest, err
}

func Base64ToZip(b64string string, zipDirectory string) error {

	// lastIdx := strings.LastIndex(zipDirectory, "/")
	// os.Mkdir(zipDirectory[:lastIdx], 0777)

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

func ZipDirectory(pathToZip string, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetPackageMetadata(directory string, metadata *models.PackageMetadata) error {
	pkgJsonPath := path.Join(directory, "package.json")
	file, err := os.Open(pkgJsonPath)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(file)
	for {
		if err := dec.Decode(metadata); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to decode JSON data from %s", pkgJsonPath)
		}
	}
	file.Close()

	(*metadata).ID = (*metadata).Name + "(" + (*metadata).Version + ")"
	//PrintMetadata(metadata)
	return err
}

func GetMetadataFromZip(zipFile string, metadata *models.PackageMetadata, readme *[]byte) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("Failed to read from %s", zipFile)
	}
	defer r.Close()
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "package.json") {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("Failed to read package.json")
			}
			dec := json.NewDecoder(rc)
			for {
				if err := dec.Decode(metadata); err == io.EOF {
					break
				} else if err != nil {
					return fmt.Errorf("Failed to parse JSON data")
				}
			}
			rc.Close()
		} else if strings.Count(f.Name, "/") == 1 {
			matched, _ := regexp.MatchString(`(?i)readme`, f.Name)
			if matched {
				//fmt.Printf("Matched: %s\n", f.Name)
				*readme = GetReadMeFromZip(f)
			}
		}
	}
	//fmt.Printf(string(readme))
	return err
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
