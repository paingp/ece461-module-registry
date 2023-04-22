package utils

import (
	"fmt"
	"tomr/models"
)

func PrintPackageData(data models.PackageData) {
	fmt.Println("Content: " + data.Content)
	fmt.Println("URL: " + data.URL)
	fmt.Println("JSProgram: " + data.JSProgram)
}

func PrintMetadata(metadata models.PackageMetadata) {
	fmt.Println("Name: " + metadata.Name)
	fmt.Println("Version: " + metadata.Version)
	fmt.Println("ID: " + metadata.ID)
	fmt.Println("License: " + metadata.License)
	fmt.Println("RepoURL: " + metadata.RepoURL)
}

func PrintRating(ratings models.PackageRating) {
	fmt.Printf("NetScore: %f\n", ratings.NetScore)
	fmt.Printf("BusFactor: %f\n", ratings.BusFactor)
	fmt.Printf("Correctness: %f\n", ratings.Correctness)
	fmt.Printf("RampUp: %f\n", ratings.RampUp)
	fmt.Printf("ResponsiveMaintainer: %f\n", ratings.ResponsiveMaintainer)
	fmt.Printf("LicenseScore: %f\n", ratings.LicenseScore)
	fmt.Printf("GoodPinningPractice: %f\n", ratings.GoodPinningPractice)
	fmt.Printf("GoodEngineeringProcess: %f\n", ratings.GoodEngineeringProcess)
}
