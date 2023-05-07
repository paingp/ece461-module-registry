package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"tomr/src/db"

	"cloud.google.com/go/storage"
)

func Regex(regex_str string) []string {
	var matches []string
	bucketName := "tomr"

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
		if err == storage.ErrObjectNotExist || err != nil {
			break
		}

		// Check match string for name, if not match check readme
		if pattern.MatchString(obj.Name) {
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
				fmt.Printf("ReadMe not found")
			}
		}

	}

	// Return list of matches
	return matches

}

type User struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}

type PackageMetadata struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	ID      string `json:"ID"`
}

type module struct {
	User     User            `json:"User"`
	Date     string          `json:"Date"`
	Metadata PackageMetadata `json:"PackageMetadata"`
	Action   string          `json:"Action"`
}

// JSON return format
// 	{
//     "User": {
//       "name": "James Davis",
//       "isAdmin": true
//     },
//     "Date": "2023-03-23T23:11:15.000Z",
//     "PackageMetadata": {
//       "Name": "Underscore",
//       "Version": "1.0.0",
//       "ID": "underscore"
//     },
//     "Action": "DOWNLOAD"
//   },

// Call function with string name of desired package to fetch history,
// username of user then string ("true" or "false") if isAdmin
func History(name string, delete int, args ...string) [][]byte {
	var user string  // Name of user
	var isAdmin bool // Is user admin or not
	// var date string		// Date of upload
	// var packName string		// Name of package
	var version string // Version of package
	// var action string	// Last action performed on a package
	var mods [][]byte

	// Setting default user parameters (user, isAdmin)
	if len(args) == 2 {
		user = args[0]
		if args[1] == "true" {
			isAdmin = true
		} else {
			isAdmin = false
		}
	} else {
		user = "ece30861defaultadminuser"
		isAdmin = true
	}

	bucketName := "tomr"

	// Create regex for string

	regex_str := "^(" + name + `)\((.*?)\)\z`
	pattern, err := regexp.Compile(regex_str)
	fmt.Print(regex_str)
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

		if isMatch && delete == 1 {
			db.DeleteObject(obj.Name)
		}

		if isMatch {
			var mod module
			rs := pattern.FindStringSubmatch(obj.Name)
			version = rs[2]
			mod.User.Name = user
			mod.User.IsAdmin = isAdmin
			mod.Metadata.Version = version
			mod.Metadata.Name = rs[1]
			mod.Metadata.ID = strings.ToLower(rs[1])
			mod.Date = "2023-03-22T23:06:25.000Z" // Default for now
			mod.Action = "CREATE"                 // Default for now
			b, err := json.MarshalIndent(mod, "", "  ")

			if err != nil {
				fmt.Println(err)
			}
			mods = append(mods, b)
		}

	}

	return mods
}

// func main() {
//     fmt.Println(History("node-du"))
// }
