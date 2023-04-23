package db

import (
	"encoding/json"
	"fmt"
	"log"
	"tomr/models"
)

const BucketName = "tomr"

func StorePackage(pkg models.PackageObject, pkgDir string) error {

	err := UploadPackage(pkgDir, pkg.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to upload package to cloud storage\n%v", err)
	}

	type objMeta struct {
		Name    string `json:"Name"`
		Version string `json:"Version"`
		ID      string `json:"ID"`
		ReadMe  string `json:"ReadMe,omitempty"`

		NetScore               string `json:"NetScore"`
		BusFactor              string `json:"BusFactor"`
		Correctness            string `json:"Correctness"`
		RampUp                 string `json:"RampUp"`
		ResponsiveMaintainer   string `json:"ResponsiveMaintainer"`
		LicenseScore           string `json:"LicenseScore"`
		GoodPinningPractice    string `json:"GoodPinningPractice"`
		GoodEngineeringProcess string `json:"GoodEngineeringProcess"`

		Content   string `json:"Content,omitempty"`
		URL       string `json:"URL,omitempty"`
		JSProgram string `json:"JSProgram,omitempty"`
	}

	objMetadata := objMeta{Name: pkg.Metadata.Name, Version: pkg.Metadata.Version, ID: pkg.Metadata.ID,
		NetScore: fmt.Sprintf("%f", pkg.Rating.NetScore), BusFactor: fmt.Sprintf("%f", pkg.Rating.BusFactor),
		Correctness: fmt.Sprintf("%f", pkg.Rating.Correctness), RampUp: fmt.Sprintf("%f", pkg.Rating.RampUp),
		ResponsiveMaintainer:   fmt.Sprintf("%f", pkg.Rating.ResponsiveMaintainer),
		LicenseScore:           fmt.Sprintf("%f", pkg.Rating.LicenseScore),
		GoodPinningPractice:    fmt.Sprintf("%f", pkg.Rating.GoodPinningPractice),
		GoodEngineeringProcess: fmt.Sprintf("%f", pkg.Rating.GoodEngineeringProcess),
		Content:                pkg.Data.Content, URL: pkg.Data.URL, JSProgram: pkg.Data.JSProgram}

	var dataMap map[string]string
	bytes, err := json.Marshal(objMetadata)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bytes, &dataMap)
	SetMetadata(dataMap, pkg.Metadata.ID)
	return err
}
