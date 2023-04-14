package ratom

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func GetToken() string {
	return os.Getenv("GITHUB_TOKEN")
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
