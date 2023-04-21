package main

import (
	"context"
	"fmt"
	"regexp"

	"cloud.google.com/go/storage"
)

func Regex(regex_str string) []string {
	var matches []string
	bucketName := "tomr-bucket"

	// Create regex for string
	pattern, err := regexp.Compile(regex_str)
	if err != nil {
		fmt.Printf("Could not create regex pattern: %v\n", err)
		return matches
	}

	// Setup GCP client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Could not setup GCP client: %v\n", err)
		return matches
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
		if err == storage.ErrObjectNotExist || err != nil{
			break
		}

		// Check match string for name, if not match check readme
		if pattern.MatchString(obj.Name) {
			fmt.Printf("Match with %s\n", obj.Name)
			matches = append(matches, obj.Name)
		} else {
			meta, err := client.Bucket(bucketName).Object(obj.Name).Attrs(ctx)
			if err != nil {
				fmt.Printf("Could not get metadata: %v\n", err)
				continue
			}
			readme, found := meta.Metadata["README"]
			if found {
				if pattern.MatchString(readme) {
					fmt.Printf("Matched in readme %s\n", obj.Name)
					matches = append(matches, obj.Name)
				}
			} else {
				fmt.Printf("readme not found")
			}
		}

	}

	// Return list of matches
	return matches

}

func main() {
	fmt.Println(Regex("(Cloudinary|lodash)"))
}
