package metrics

import (
	"context"
	"log"
	"os"
	"tomr/models"
	"tomr/src/utils"

	"golang.org/x/oauth2"
)

func RatePackage(url string, pkgDirectory string, rating *models.PackageRating, license string, readMe *[]byte) error {
	gitEndpoint := utils.GetGithubEndpoint(url)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	jsonData, err := utils.GetDataFromGithub(httpClient, gitEndpoint)
	if err != nil {
		log.Fatalf("Failed to get data from GITHUB API rate package with URL: %s\n", url)
	}

	(*rating).BusFactor = getBusFactor(jsonData)
	(*rating).Correctness = getCorrectness(jsonData)
	(*rating).RampUp = getRampUp(jsonData, httpClient)
	(*rating).ResponsiveMaintainer = getResponsiveMaintainer(jsonData)

	(*rating).LicenseScore = getLicenseScore(license, pkgDirectory, readMe)
	(*rating).GoodPinningPractice = getGoodPinningPractices(gitEndpoint, httpClient)
	(*rating).GoodEngineeringProcess = getGoodEngineeringProcess(gitEndpoint, httpClient, pkgDirectory)
	(*rating).NetScore = getNetScore(*rating)

	//os.RemoveAll("src/metrics/temp")
	//utils.PrintRating((*rating))

	return err
}
