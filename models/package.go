package models

// This is a \"union\" type. - On package upload, either Content or URL should be set. - On package update, exactly one field should be set. - On download, the Content field should be set.
type PackageData struct {
	// Package contents. This is the zip file uploaded by the user. (Encoded as text using a Base64 encoding).  This will be a zipped version of an npm package's GitHub repository, minus the \".git/\" directory.\" It will, for example, include the \"package.json\" file that can be used to retrieve the project homepage.  See https://docs.npmjs.com/cli/v7/configuring-npm/package-json#homepage.
	Content string `json:"Content,omitempty"`
	// Package URL (for use in public ingest).
	URL string `json:"URL,omitempty"`
	// A JavaScript program (for use with sensitive modules).
	JSProgram string `json:"JSProgram,omitempty"`
}

// The \"Name\" and \"Version\" are used as a unique identifier pair when uploading a package.  The \"ID\" is used as an internal identifier for interacting with existing packages.
type PackageMetadata struct {
	Name string `json:"Name"`
	// Package version
	Version string `json:"Version"`

	ID string `json:"ID"`

	License string `json:"License"`

	RepoURL string `json:"Homepage"`

	Date string `json:"Date"`

	Action string `json:"Action"`
}

type PackageObject struct {
	Metadata *PackageMetadata `json:"metadata"`

	Data *PackageData `json:"data"`

	Rating *PackageRating `json:"rating"`
}

// Package rating (cf. Project 1).  If the Project 1 that you inherited does not support one or more of the original properties, denote this with the value \"-1\".
type PackageRating struct {
	NetScore float64 `json:"NetScore"`

	BusFactor float64 `json:"BusFactor"`

	Correctness float64 `json:"Correctness"`

	RampUp float64 `json:"RampUp"`

	ResponsiveMaintainer float64 `json:"ResponsiveMaintainer"`

	LicenseScore float64 `json:"LicenseScore"`
	// The fraction of its dependencies that are pinned to at least a specific major+minor version, e.g. version 2.3.X of a package. (If there are zero dependencies, they should receive a 1.0 rating. If there are two dependencies, one pinned to this degree, then they should receive a Â½ = 0.5 rating).
	GoodPinningPractice float64 `json:"GoodPinningPractice"`
	// The fraction of project code that was introduced hrough pull requests with a code review).
	GoodEngineeringProcess float64 `json:"GoodEngineeringProcess"`
}

type PackageQuery struct {
	Version string `json:"Version,omitempty"`

	Name string `json:"Name"`
}
