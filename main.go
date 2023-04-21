package main

import (
	"fmt"
	"os"
	"tomr/src/metrics"
)

func main() {

	url := "https://www.npmjs.com/package/du"

	fmt.Println(os.Getenv("GITHUB_TOKEN"))

	metrics.RatePackage(url)

	//dir := "C:/Users/paing/Desktop/College Assignments/Spring 2023/ECE 30861/dev/server/temp/package.zip"

	//handlers.GetPackageMetadata(dir)
}
