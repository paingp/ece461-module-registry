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

	(*rating).BusFactor = -1.0
	(*rating).Correctness = -1.0
	(*rating).RampUp = -1.0
	(*rating).ResponsiveMaintainer = -1.0
	(*rating).LicenseScore = 0.0

	(*rating).GoodPinningPractice = -1.0
	(*rating).GoodEngineeringProcess = -1.0

	(*rating).BusFactor = getBusFactor(jsonData)
	if ingestion {
		if (*rating).BusFactor < 0.5 {
			fmt.Printf("BusFactor score of %f doesn't meet criteria for ingestion", (*rating).BusFactor)
			return fmt.Errorf("BusFactor score of %f doesn't meet criteria for ingestion", (*rating).BusFactor)
		}
	}

	(*rating).Correctness = getCorrectness(jsonData)

	if ingestion {
		if (*rating).Correctness < 0.5 {
			fmt.Printf("Correctness score of %f  doesn't meet criteria for ingestion", (*rating).Correctness)
			return fmt.Errorf("Correctness score of %f  doesn't meet criteria for ingestion", (*rating).Correctness)
		}
	}

	(*rating).RampUp = getRampUp(jsonData, httpClient)
	if ingestion {
		if (*rating).RampUp < 0.5 {
			fmt.Printf("RampUp score of %f doesn't meet for ingestion", (*rating).RampUp)
			return fmt.Errorf("RampUp score of %f doesn't meet for ingestion", (*rating).RampUp)
		}
	}

	(*rating).ResponsiveMaintainer = getResponsiveMaintainer(jsonData)
	if ingestion {
		if (*rating).ResponsiveMaintainer < 0.5 {
			fmt.Printf("ResponsiveMaintainer score of %f doesn't meet criteria for ingestion", (*rating).ResponsiveMaintainer)
			return fmt.Errorf("ResponsiveMaintainer score of %f doesn't meet criteria for ingestion", (*rating).ResponsiveMaintainer)
		}
	}

	(*rating).LicenseScore = getLicenseScore(license, pkgDirectory, readMe)
	if ingestion {
		if (*rating).LicenseScore < 0.5 {
			fmt.Printf("License does not meet criteria for ingestion (must be ompatible with LGPLv2.1)")
			return fmt.Errorf("License does not meet criteria for ingestion (must be ompatible with LGPLv2.1)")
		}
	}

	(*rating).GoodPinningPractice = getGoodPinningPractices(gitEndpoint, httpClient)
	(*rating).GoodEngineeringProcess = getGoodEngineeringProcess(gitEndpoint, httpClient, pkgDirectory)
	// (*rating).GoodEngineeringProcess = 0.0

	(*rating).NetScore = getNetScore(*rating)

	//os.RemoveAll("src/metrics/temp")
	//utils.PrintRating((*rating))

	return err
}
