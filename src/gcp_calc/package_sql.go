package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/go-sql-driver/mysql"
)
var db *sql.DB

type Metadata struct {
    ID     					string
	NAME					string
	Version					string
	Licesnse				string
	ReadMe					string
	RepoURL					string
	NetScore				float64
	BusFactor				float64
	Correctness				float64
	RampUp					float64
	ResponsiveMaintainer	float64
	LicenseScore			float64
	GoodPinningPractice		float64
	GoodEngineeringPractice	float64
}

func AddFileMetadata(id string, name string, version string, license string, readme string, repo string) {
    // Capture connection properties.
    cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "34.122.232.231",
        DBName: "File_Metadata",
    }
    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

	var meta Metadata

	meta.NAME = "temp"
	meta.ID = "temp"
	meta.Version = "temp"

	addMetadata(meta)
}

func addMetadata(meta Metadata) (int64, error) {
    result, err := db.Exec("INSERT INTO PackageMetadata (ID, NAME, Version) VALUES (?, ?, ?)", meta.ID, meta.NAME, meta.Version)
    if err != nil {
        return 0, fmt.Errorf("addAlbum: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addAlbum: %v", err)
    }
    return id, nil
}

func main() {
	AddFileMetadata("hello", "hello", "hello", "hello", "hello", "hello")
}