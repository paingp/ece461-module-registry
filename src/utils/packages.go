package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
)

type pack struct {
	Version string `json:"Version"`
	Name    string `json:"Name"`
}

// Pass in as version, name
func Packages(version string, name string) [][]byte {

	var regex_str string

	var packs [][]byte

	bucketName := "tomr"

	// Create regex for string

	if name != "*" {
		regex_str = "^(" + name + `)\((.*?)\)\z`
	} else {
		regex_str = "^(" + ".*?" + `)\((.*?)\)\z`
	}
	pattern, err := regexp.Compile(regex_str)
	if err != nil {
		fmt.Printf("Could not create regex pattern: %v\n", err)
		return nil
	}

	// Setup GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Could not setup GCP client: %v\n", err)
		return nil
	}
	defer client.Close()

	fmt.Print("are we here")

	// Create a query
	query := &storage.Query{
		Delimiter: "/",
	}

	// Get list of all objects and iterate
	blob := client.Bucket(bucketName).Objects(ctx, query)
	for {
		// Get next object and check for end of bucket
		obj, err := blob.Next()
		if err == storage.ErrObjectNotExist || err != nil {
			break
		}

		isMatch := false

		// Check match string for name, if not match check readme
		if pattern.MatchString(obj.Name) {
			// fmt.Printf("Match with %s\n", obj.Name)
			isMatch = true
		}

		if isMatch {
			rs := pattern.FindStringSubmatch(obj.Name)
			if !strings.Contains(version, rs[2]) {
				continue
			}

			var pack pack

			pack.Version = rs[2]
			pack.Name = obj.Name

			b, err := json.MarshalIndent(pack, "  ", "  ")

			if err != nil {
				fmt.Println(err)
			}

			packs = append(packs, b)
		}

	}

	return packs
}

const pkgDirPath = "src/metrics/temp"

// func CreatePackage(content string, url string, jsprogram string) {

// 	packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
// 	pkgDir := ""
// 	metadata := models.PackageMetadata{}
// 	//utils.PrintPackageData(packageData)
// 	rating := models.PackageRating{}
// 	var readMe []byte
// 	// Return Error 400 if both Content and URL are set
// 	if (packageData.Content != "") && (packageData.URL != "") {
// 		fmt.Printf("Error 400: Content and URL cannot be both set")
// 	} else if packageData.Content != "" { // Only Content is set
// 		// Decode base64 string into zip
// 		pkgDir = path.Join(pkgDirPath, "package.zip")
// 		Base64ToZip(packageData.Content, pkgDir)
// 		err := GetMetadataFromZip(pkgDir, &metadata, &readMe)
// 		if err != nil {
// 			log.Fatalf("Failed to get metadata from zip file\n")
// 		}
// 		metadata.ID = metadata.Name + "(" + metadata.Version + ")"
// 		err = metrics.RatePackage(metadata.RepoURL, pkgDir, &rating, metadata.License, &readMe)
// 		if err != nil {
// 			log.Fatalf("Failed to get metadata from zip file\n")
// 		}
// 	} else { // Only URL is set
// 		gitUrl := GetGithubUrl(url)
// 		pkgDir = CloneRepo(gitUrl, pkgDirPath)
// 		err := metrics.RatePackage(gitUrl, pkgDir, &rating, "", nil)
// 		if err != nil {
// 			log.Fatalf("Failed to rate package at URL: %s\n", url)
// 		}
// 		// Check if package meets criteria for ingestion
// 		GetPackageMetadata(pkgDir, &metadata)
// 		err = ZipDirectory(pkgDir, pkgDir+".zip")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		pkgDir += ".zip"
// 	}
// 	PrintMetadata(metadata)
// 	PrintRating(rating)
// 	// Upload package and store data in system
// 	pkg := models.PackageObject{Metadata: &metadata, Data: &packageData, Rating: &rating}

// 	err := db.StorePackage(pkg, pkgDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
// 	packages("Exact (1.0.0)\nBounded range (1.2.3-2.1.0)\nCarat (^1.9.4)\nTilde (~1.2.0)", "du")
// }
