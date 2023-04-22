package main

import (
	"log"
	"tomr/src/handlers"
	"tomr/src/utils"
)

func main() {

	//fmt.Println(os.Getenv("GITHUB_TOKEN"))

	//url := "https://www.npmjs.com/package/du"
	//gitUrl := utils.GetGithubUrl(url)
	//utils.CloneRepo(gitUrl, "src/metrics/temp")

	/*
		pkgDir := "src/metrics/temp/package.zip"
		metadata := models.PackageMetadata{}
		var readme []byte
		utils.GetMetadataFromZip(pkgDir, &metadata, &readme)
		fmt.Printf(string(readme))
	*/
	//pkgDir := ""

	// Encode/decode between base64 string and ZIP

	content, err := utils.ZipToBase64("src/metrics/temp/node-du.zip")
	if err != nil {
		log.Fatal(err)
	}

	/*
		err = utils.Base64ToZip(content, "src/metrics/temp/package.zip")
		if err != nil {
			log.Fatal(err)
		}
	*/
	handlers.CreatePackage(content, "", "console.log('Hello World')")

	//metrics.RatePackage(url, pkgDir)

}
