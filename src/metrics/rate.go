package metrics

import (
	"context"
	"fmt"
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
		fmt.Print("here")
		log.Fatalf("Failed to get data from GITHUB API rate package with URL: %s\n", url)
	}

	ingestion := false
	if readMe == nil {
		ingestion = true
	}

	(*rating).BusFactor = getBusFactor(jsonData)
	if ingestion {
		if (*rating).BusFactor < 0.5 {
			return fmt.Errorf("score too low for BusFactor to meet criteria for ingestion")
		}
	}
	(*rating).Correctness = getCorrectness(jsonData)
	if ingestion {
		if (*rating).Correctness < 0.5 {
			return fmt.Errorf("score too low for correctness to meet criteria for ingestion")
		}
	}
	(*rating).RampUp = getRampUp(jsonData, httpClient)
	// fmt.Print("Ramp up Score", (*rating).RampUp)
	if ingestion {
		if (*rating).RampUp < 0.5 {
			return fmt.Errorf("score too low for rampup to meet criteria for ingestion")
		}
	}
	(*rating).ResponsiveMaintainer = getResponsiveMaintainer(jsonData)
	if ingestion {
		if (*rating).ResponsiveMaintainer < 0.5 {
			return fmt.Errorf("score too low for responsive maintainer to meet criteria for ingestion")
		}
	}
	if readMe == nil {
		(*rating).LicenseScore = getLicenseScore(license, pkgDirectory, nil)
	} else {
		(*rating).LicenseScore = getLicenseScore(license, pkgDirectory, *readMe)
	}

	if ingestion {
		if (*rating).LicenseScore < 0.5 {
			return fmt.Errorf("score too low for license to meet criteria for ingestion")
		}
	}

	(*rating).GoodPinningPractice = getGoodPinningPractices(gitEndpoint, httpClient)
	// (*rating).GoodEngineeringProcess = getGoodEngineeringProcess(gitEndpoint, httpClient, pkgDirectory)
	(*rating).GoodEngineeringProcess = 0.0

	(*rating).NetScore = getNetScore(*rating)

	//os.RemoveAll("src/metrics/temp")
	utils.PrintRating((*rating))

	return err
}
