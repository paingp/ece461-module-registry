package metrics

import (
	"context"
	"log"
	"os"
	"tomr/models"
	"tomr/src/utils"

	"golang.org/x/oauth2"
)

func RatePackage(url string, pkgDirectory string, license string, readMe *[]byte) error {
	//gitEndpoint := utils.GetGithubEndpoint(url)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	jsonData, err := utils.GetDataFromGithub(httpClient, url)
	if err != nil {
		log.Fatalf("Failed to get data from GITHUB API rate package with URL: %s\n", url)
	}
	rating := models.PackageRating{}

	rating.BusFactor = getBusFactor(jsonData)
	rating.Correctness = getCorrectness(jsonData)
	rating.RampUp = getRampUp(jsonData, httpClient)
	rating.ResponsiveMaintainer = getResponsiveMaintainer(jsonData)

	rating.LicenseScore = getLicenseScore(license, pkgDirectory, readMe)
	rating.GoodPinningPractice = getGoodPinningPractices(url, httpClient)
	//rating.GoodEngineeringProcess = getGoodEngineeringProcess(url, httpClient, pkgDirectory)

	//os.RemoveAll("src/metrics/temp")
	utils.PrintRating(rating)

	return nil
}
