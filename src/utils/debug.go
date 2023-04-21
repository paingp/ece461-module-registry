package utils

import (
	"fmt"
	"tomr/models"
)

func PrintMetadata(metadata models.PackageMetadata) {
	fmt.Println("Name: " + metadata.Name)
	fmt.Println("Version: " + metadata.Version)
	fmt.Println("ID: " + metadata.ID)
	fmt.Println("License: " + metadata.License)
	fmt.Println("RepoURL: " + metadata.RepoURL)
}
