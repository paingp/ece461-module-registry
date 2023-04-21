package metrics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"tomr/src/utils"

	"golang.org/x/oauth2"
)

func RatePackage(url string) {
	gitUrl := utils.GetGithubUrl(url)
	gitEndpoint := utils.GetGithubEndpoint(gitUrl)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	resp, error := httpClient.Get(gitEndpoint)

	if (error != nil) || (resp.StatusCode != http.StatusOK) {
		fmt.Println(resp.StatusCode)
	}

	/*
		jsonData := utils.GetDataFromGithub(httpClient, gitEndpoint)

		rating := models.PackageRating{}
		rating.BusFactor = getBusFactor(jsonData)
		rating.Correctness = getCorrectness(jsonData)
		rating.RampUp = getRampUp(jsonData, httpClient)
		rating.ResponsiveMaintainer = getResponsiveMaintainer(jsonData)
		//rating.LicenseScore = getLicenseScore()
		rating.GoodPinningPractice = getGoodPinningPractices(gitEndpoint, httpClient)
		//rating.GoodEngineeringProcess = getGoodEngineeringProcess(jsonData, )
	*/
}
