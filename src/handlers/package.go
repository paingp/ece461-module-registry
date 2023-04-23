package handlers

import (
	"fmt"
	"log"
	"path"
	"tomr/models"
	"tomr/src/db"
	"tomr/src/metrics"
	"tomr/src/utils"
)

const pkgDirPath = "src/metrics/temp" // temp directory to store packages

func CreatePackage(content string, url string, jsprogram string) {
	packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
	pkgDir := ""
	metadata := models.PackageMetadata{}
	//utils.PrintPackageData(packageData)
	rating := models.PackageRating{}
	var readMe []byte
	// Return Error 400 if both Content and URL are set
	if (packageData.Content != "") && (packageData.URL != "") {
		fmt.Printf("Error 400: Content and URL cannot be both set")
	} else if packageData.Content != "" { // Only Content is set
		// Decode base64 string into zip
		pkgDir = path.Join(pkgDirPath, "package.zip")
		utils.Base64ToZip(packageData.Content, pkgDir)
		err := utils.GetMetadataFromZip(pkgDir, &metadata, &readMe)
		if err != nil {
			log.Fatalf("Failed to get metadata from zip file\n")
		}
		metadata.ID = metadata.Name + "(" + metadata.Version + ")"
		err = metrics.RatePackage(metadata.RepoURL, pkgDir, &rating, metadata.License, &readMe)
	} else { // Only URL is set
		gitUrl := utils.GetGithubUrl(url)
		pkgDir = utils.CloneRepo(gitUrl, pkgDirPath)
		err := metrics.RatePackage(gitUrl, pkgDir, &rating, "", nil)
		if err != nil {
			log.Fatalf("Failed to rate package at URL: %s\n", url)
		}
		// Check if package meets criteria for ingestion
		utils.GetPackageMetadata(pkgDir, &metadata)
		err = utils.ZipDirectory(pkgDir, pkgDir+".zip")
		if err != nil {
			log.Fatal(err)
		}
		pkgDir += ".zip"
	}
	utils.PrintMetadata(metadata)
	utils.PrintRating(rating)
	fmt.Printf(pkgDir)
	// Upload package and store data in system
	pkg := models.PackageObject{Metadata: &metadata, Data: &packageData, Rating: &rating}

	err := db.StorePackage(pkg, pkgDir)
	if err != nil {
		log.Fatal(err)
	}
}
