package utils

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"encoding/json"

)

type pack struct {
	Version string `json:"Version"`
	Name 	string `json:"Name"`
}


// Pass in as version, name
func packages(version string, name string) [][]byte{

	var regex_str string

	var packs [][]byte

	bucketName := "tomr"

	// Create regex for string
	
	if name != "*" {
		regex_str = "^(" + name + `)\((.*?)\)\z`
	} else {
		regex_str = "^(" + ".*?" + `)\((.*?)\)\z`
	}
	pattern, err := regexp.Compile(regex_str)
	if err != nil {
		fmt.Printf("Could not create regex pattern: %v\n", err)
		return nil
	}

	// Setup GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Could not setup GCP client: %v\n", err)
		return nil
	}
	defer client.Close()

	// Create a query
	query := &storage.Query{
		Delimiter: "/",
	}

	// Get list of all objects and iterate
	blob := client.Bucket(bucketName).Objects(ctx, query)
	for {
		// Get next object and check for end of bucket
		obj, err := blob.Next()
		if err == storage.ErrObjectNotExist || err != nil {
			break
		}

		isMatch := false

		// Check match string for name, if not match check readme
		if pattern.MatchString(obj.Name) {
			// fmt.Printf("Match with %s\n", obj.Name)
			isMatch = true
		} 

		if isMatch {
			rs := pattern.FindStringSubmatch(obj.Name)
			if(!strings.Contains(version, rs[2])) {
				continue
			}
			
			var pack pack

			pack.Version = rs[2]
			pack.Name = obj.Name

			b, err := json.MarshalIndent(pack, "  ", "  ")

			if err != nil {
				fmt.Println(err)
			}

			packs = append(packs, b)
		}

	}

	return packs
}


// func main() {
// 	packages("Exact (1.0.0)\nBounded range (1.2.3-2.1.0)\nCarat (^1.9.4)\nTilde (~1.2.0)", "du")
// }
