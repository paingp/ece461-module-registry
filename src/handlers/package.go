package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"tomr/models"
	"tomr/src/db"
	"tomr/src/metrics"
	"tomr/src/utils"

	"github.com/go-chi/chi/v5"
)

const PkgDirPath = "src/metrics/temp" // temp directory to store packages
const auth_success = "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiZWNlMzA4NjFkZWZhdWx0YWRtaW51c2VyIiwicGFzc3dvcmQiOiJjb3JyZWN0aG9yc2ViYXR0ZXJ5c3RhcGxlMTIzKCFfXytAKiooQeKAmeKAnWA7RFJPUCBUQUJMRSBwYWNrYWdlczsifQ.TSGs6VJMFx5NV2RoHrhEP_FK8nv4Wlzc4gQls2JYPC4"
const BucketName = "tomr"

type metadata struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	ID      string `json:"ID"`
}

type Package struct {
	Metadata metadata           `json:"metadata"`
	Data     models.PackageData `json:"data"`
}

func CreatePackage(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering CreatePackage()\n")

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	if given_xAuth == auth_success {
		var data models.PackageData
		body, _ := ioutil.ReadAll(request.Body)
		json.Unmarshal(body, &data)

		content := data.Content
		url := data.URL
		fmt.Print("URL: " + url + "\n")
		jsprogram := data.JSProgram

		packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
		pkgDir := ""
		metadata := models.PackageMetadata{}
		//utils.PrintPackageData(packageData)
		rating := models.PackageRating{}
		var readMe []byte

		// Return Error 400 if both Content and URL are set
		if (packageData.Content != "") && (packageData.URL != "") {

			fmt.Print("Exiting CreatePackage() due to missing field\n")

			writer.WriteHeader(400)
			writer.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
			os.RemoveAll(PkgDirPath)
			return

		} else if packageData.Content != "" { // Only Content is set
			// Decode base64 string into zip
			pkgDir = path.Join(PkgDirPath, "package.zip")
			utils.Base64ToZip(packageData.Content, pkgDir)
			err := utils.GetMetadataFromZip(pkgDir, &metadata, &readMe)
			if err != nil {
				writer.WriteHeader(400)
				writer.Write([]byte("Error no metadata in zip"))
				// log.Fatalf("Failed to get metadata from zip file\n")
				os.RemoveAll(PkgDirPath)
				return
			}
			metadata.ID = metadata.Name + "(" + metadata.Version + ")"

			fmt.Print("Content (not URL) found for " + string(metadata.ID) +" in CreatePackage()\n")

			if db.DoesPackageExist(metadata.ID) {
				writer.WriteHeader(409)
				writer.Write([]byte("Package exists already."))
				os.RemoveAll(PkgDirPath)
				return
			}

			err = metrics.RatePackage(metadata.RepoURL, pkgDir, &rating, metadata.License, &readMe)
			if err != nil {
				writer.WriteHeader(424)
				writer.Write([]byte("Package is not uploaded due to the disqualified rating."))
				
				log.Fatalf("Failed to get metadata from zip file in CreatePackage()\n")
				
				os.RemoveAll(PkgDirPath)
				return
			}

			writer.WriteHeader(201)
			writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))

		} else { // Only URL is set
			gitUrl := utils.GetGithubUrl(url)
			pkgDir = utils.CloneRepo(gitUrl, PkgDirPath)
			err := metrics.RatePackage(gitUrl, pkgDir, &rating, "", nil)
			if err != nil {

				writer.WriteHeader(424)
				writer.Write([]byte("Package is not uploaded due to the disqualified rating."))
				// log.Fatalf("Failed to get metadata from zip file\n")
				os.RemoveAll(PkgDirPath)
				return
			}

			// Check if package meets criteria for ingestion
			utils.GetPackageMetadata(pkgDir, &metadata)
			
			fmt.Print("URL (not Content) found for " + string(metadata.ID) +" in CreatePackage()\n")

			if db.DoesPackageExist(metadata.ID) {
				writer.WriteHeader(409)
				writer.Write([]byte("Package exists already."))
				os.RemoveAll(PkgDirPath)
				return
			}

			err = utils.ZipDirectory(pkgDir, pkgDir+".zip")
			if err != nil {
				writer.WriteHeader(400)
				writer.Write([]byte("Unable to zip directory"))
				// log.Fatalf("Failed to get metadata from zip file\n")
				os.RemoveAll(PkgDirPath)
				return
			}
			pkgDir += ".zip"

			writer.WriteHeader(201)
			writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID.\n"))
		}

		fmt.Print("Exiting if-else if condition in CreatePackage()\n")

		pkg := models.PackageObject{Metadata: &metadata, Data: &packageData, Rating: &rating}

		base64, err := utils.ZipToBase64(pkgDir)

		if err != nil {
			os.RemoveAll(PkgDirPath)
			fmt.Print("could not convert to base64")
			return
		}

		writePath := "src/db/upload.txt"
		fileWriter, err := os.Create(writePath)
		fileWriter.WriteString(base64)

		if err != nil {
			os.RemoveAll(writePath)
			os.RemoveAll(PkgDirPath)
			fmt.Print("Failed to write base64 encoding : ", err)
			return
		}

		return_json, err := db.StorePackage(pkg, writePath)
		if err != nil {
			os.RemoveAll(writePath)
			os.RemoveAll(PkgDirPath)
			fmt.Print("did not get return json")
		}

		var returnVal db.Return_storage
		json.Unmarshal(return_json, &returnVal)

		returnVal.Data.Content = base64
		return_json, _ = json.MarshalIndent(returnVal, "", "  ")

		writer.Write([]byte(string(return_json)))

		os.RemoveAll(writePath)
		os.RemoveAll(PkgDirPath)

		fmt.Print("Exiting CreatePackage() successfully after removing temp directories\n")

	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func GetPackageByRegEx(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering GetPackageByRegex()\n")

	request.ParseForm()

	var given_xAuth string
	var regex_string string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
		regex_string = request.Form["Regex String"][0]
	} else {
		given_xAuth = request.Header["X-Authorization"][0]

		type regex_body struct {
			Regex string `json:"Regex"`
		}

		body, _ := io.ReadAll(request.Body)

		var regex regex_body
		json.Unmarshal(body, &regex)

		regex_string = string(body)
	}

	fmt.Print("Recieved RegEx string correctly in GetPackageByRegex()\n")

	if given_xAuth == auth_success && regex_string != "" {

		return_string := utils.Regex(regex_string)

		if len(return_string) == 0 {
			writer.WriteHeader(404)
			writer.Write([]byte("Status: 404 No package found under this regex."))
			return
		}

		type Regex_output struct {
			Version string
			Name    string
		}

		regex_str := "(.*?)" + `\((.*?)\)`
		pattern, _ := regexp.Compile(regex_str)

		var regex_return []Regex_output

		for _, i := range return_string {
			rs := pattern.FindStringSubmatch(i)

			regex_singleMatch := Regex_output{rs[2], rs[1]}
			regex_return = append(regex_return, regex_singleMatch)
		}

		writer.WriteHeader(200)
		writer.Write([]byte("["))

		for i := 0; i < len(regex_return); i++ {
			output_string := "\n  {\n    \"Version\": \"" + regex_return[i].Version + "\",\n    \"Name\": \"" + regex_return[i].Name + "\"\n  }"
			writer.Write([]byte(output_string))

			if i == len(regex_return)-1 {
				writer.Write([]byte(","))
			}
		}
		writer.Write([]byte("\n]"))

		fmt.Print("Matches may have been found exiting GetPackageByRegex() correctly\n")

	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("Status: 400 There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func ResetRegistry(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering ResetRegistry()\n")

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
		// given_xAuth = chi.URLParam(request, "auth")
	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	fmt.Print("Passed authentication in ResetRegistry()\n")

	if given_xAuth == auth_success {

		// cmd := exec.Command("python3", "src/gcp_calc/deleteBucket.py", "tomr")
		// err := cmd.Run()

		// if err != nil {
		// 	fmt.Print("Error here\n")
		// 	println(err.Error())
		// 	return
		// }

		// output, err := cmd.CombinedOutput()
		// if err != nil {
		// 	fmt.Println(fmt.Sprint(err) + ": " + string(output))
		// 	return
		// }

		fmt.Print("Calling db.DeleteBucketContets() in ResetRegistry()\n")

		err := db.DeleteBucketContents()

		if err != nil {
			log.Fatal("Could not delete bucket contents in ResetRegistry()\n")
		}

		writer.WriteHeader(200)
		writer.Write([]byte("Registry is reset."))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

	fmt.Print("Exiting ResetRegistry() correctly\n")

}

func RetrievePackage(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering RetrievePackage()\n")

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	id := chi.URLParam(request, "id")

	fmt.Print("Retrieving package of " + id + " in RetrievePackage()\n")

	if given_xAuth == auth_success {

		filepath := "src/handlers/readTo.txt"

		errFile := db.DownloadFile(id, filepath)

		if errFile != nil {
			writer.WriteHeader(404)
			writer.Write([]byte("Package does not exist."))
			return
		}

		dat, err := os.ReadFile(filepath)

		if err != nil {
			panic(err)
		}

		base64 := string(dat) // get file contents

		regex_str := "(.*?)" + `\((.*?)\)`
		pattern, _ := regexp.Compile(regex_str)
		rs := pattern.FindStringSubmatch(id)

		// version rs[2]
		// name rs[1]

		attrs, _ := db.GetMetadata(BucketName, id)
		metadata := attrs.Metadata

		var return_package Package
		return_package.Metadata.Name = rs[1]
		return_package.Metadata.Version = rs[2]
		return_package.Metadata.ID = metadata["ID"]
		return_package.Data.Content = base64
		return_package.Data.URL = metadata["URL"]
		return_package.Data.JSProgram = metadata["JSProgram"]

		fmt.Print("Correctly getting information in RetrievePackage()\n")

		b, err := json.MarshalIndent(return_package, "", "  ")

		if err != nil {
			fmt.Println(err)
		}

		writer.WriteHeader(200)
		writer.Write([]byte(string(b)))
	} else if given_xAuth == "" || id == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.Write([]byte("{\n  \"code\": 0,\n  \"message\": \"Other Error\"\n}"))
	}

	fmt.Print("Properly exiting RetrievePackage()\n")

}

func UpdatePackage(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering UpdatePackage()\n")

	request.ParseForm()

	var given_xAuth string
	var id string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	id = chi.URLParam(request, "id")

	if given_xAuth == auth_success {

		if !db.DoesPackageExist(id) {
			writer.WriteHeader(404)
			writer.Write([]byte("Package does not exist"))
			return
		}

		fmt.Print("Deleting prior version in RetrievePackage()\n")

		db.DeleteObject(id)

		var recieve_package Package
		body, _ := io.ReadAll(request.Body)
		json.Unmarshal(body, &recieve_package)

		content := recieve_package.Data.Content
		//url := recieve_package.Data.URL
		//jsprogram := recieve_package.Data.JSProgram

		attrs, err := db.GetMetadata(BucketName, id)
		if err != nil {
			fmt.Println(err)
		}
		objMetadata := attrs.Metadata

		writePath := "src/db/upload.txt"
		fileWriter, err := os.Create(writePath)
		fileWriter.WriteString(content)

		if err != nil {
			os.RemoveAll(writePath)
			fmt.Printf("Failed to write base64 encoding : %s\n", err)
			return
		}

		// db.DeleteObject(id)

		fmt.Print("Uploading new version in RetrievePackage()\n")
		err = db.UploadPackage(writePath, id)
		if err != nil {
			fmt.Printf("Failed to upload package to cloud storage\n%v", err)
		}

		fmt.Print("Setting metadata of new version in RetrievePackage()\n")
		db.SetMetadata(objMetadata, id)

		writer.WriteHeader(200)
		writer.Write([]byte("Version is updated"))

	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

	fmt.Print("Properly exitting RetrievePackage()\n")
}

func DeletePackage(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering DeletePackage()\n")

	request.ParseForm()

	var given_xAuth string
	var id string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	id = chi.URLParam(request, "id")

	fmt.Print("Deleting " + id + " DeletePackage()\n")

	if given_xAuth == auth_success {

		if !db.DoesPackageExist(id) {
			writer.WriteHeader(404)
			writer.Write([]byte("Package does not exist"))
			return
		}

		db.DeleteObject(id)
		fmt.Print("Package deleted DeletePackage()\n")
    
		writer.WriteHeader(200)
		writer.Write([]byte("Version is deleted."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

	fmt.Print("Properly exitting DeletePackage()\n")
}

func ListPackages(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering ListPackages()\n")

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	if given_xAuth == auth_success {

		type Pack struct {
			Version string `json:"Version"`
			Name    string `json:"Name"`
		}

		var Packs []Pack
		c1 := make(chan []Pack, 1)

		go func() {
			err := json.NewDecoder(request.Body).Decode(&Packs)
			if err != nil {
				return
			}
			c1 <- Packs
		}()

		select {
		case res := <-c1:
			writer.WriteHeader(200)

			writer.Write([]byte("[\n"))

			fmt.Print("Searching existing packages in ListPackages()\n")

			for i := 0; i < len(res); i++ {
				result := utils.Packages(res[i].Version, res[i].Name)
				writer.Write([]byte("  "))

				writer.Write([]byte(string(result[0])))

				if i != len(res)-1 {
					writer.Write([]byte(",\n"))
				}
			}

			writer.Write([]byte("\n]"))
		case <-time.After(60 * time.Second):

			fmt.Print("ERROR: Searching existing packages in ListPackages() timed out\n")

			writer.WriteHeader(413)
			writer.Write([]byte("Too many packages returned."))
		}

	} else if given_xAuth == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.Write([]byte("{\n  \"code\": 0,\n  \"message\": \"Other Error\"\n}"))
	}

	fmt.Print("Properly exitting ListPackages()\n")
}

func RatePackage(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering RatePackage()\n")

	request.ParseForm()

	var given_xAuth string
	var id string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	id = chi.URLParam(request, "id")

	fmt.Print("Rating package: " + id + " in RatePackage()\n")

	if given_xAuth == auth_success {

		if !db.DoesPackageExist(id) {
			writer.WriteHeader(404)
			writer.Write([]byte("Package does not exist."))
			return
		}

		type ratings struct {
			GoodPinningPractice  string
			NetScore             string
			PullRequest          string
			ResponsiveMaintainer string
			LicenseScore         string
			RampUp               string
			BusFactor            string
			Correctness          string
		}

		attrs, _ := db.GetMetadata(BucketName, id)
		metadata := attrs.Metadata

		var package_ratings ratings

		package_ratings.GoodPinningPractice = metadata["GoodPinningPractice"]
		tempFloat, err := strconv.ParseFloat(package_ratings.GoodPinningPractice, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.NetScore = metadata["NetScore"]
		tempFloat, err = strconv.ParseFloat(package_ratings.NetScore, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.PullRequest = metadata["GoodEngineeringProcess"]
		tempFloat, err = strconv.ParseFloat(package_ratings.PullRequest, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.ResponsiveMaintainer = metadata["ResponsiveMaintainer"]
		tempFloat, err = strconv.ParseFloat(package_ratings.ResponsiveMaintainer, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.LicenseScore = metadata["LicenseScore"]
		tempFloat, err = strconv.ParseFloat(package_ratings.LicenseScore, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.RampUp = metadata["RampUp"]
		tempFloat, err = strconv.ParseFloat(package_ratings.RampUp, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.BusFactor = metadata["BusFactor"]
		tempFloat, err = strconv.ParseFloat(package_ratings.BusFactor, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		package_ratings.Correctness = metadata["Correctness"]
		tempFloat, err = strconv.ParseFloat(package_ratings.Correctness, 32)
		if err != nil || tempFloat < 0 || tempFloat > 1 {
			writer.WriteHeader(500)
			writer.Write([]byte("The package rating system choked on at least one of the metrics."))
		}

		fmt.Print("All scores calculated in RatePackage()\n")

		writer.WriteHeader(200)
		return_json, _ := json.MarshalIndent(package_ratings, "", "  ")

		writer.Write([]byte(string(return_json)))
		// writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

	fmt.Print("Exitting RatePackage()\n")
}

func CreateAuthToken(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering CreateAuthToken()\n")

	type User_struct struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"isAdmin"`
	}

	type Secret_struct struct {
		Password string `json:"password"`
	}

	type Auth struct {
		User   User_struct   `json:"user"`
		Secret Secret_struct `json:"secret"`
	}

	var auth_struct Auth
	body, _ := io.ReadAll(request.Body)
	// fmt.Print(body)
	json.Unmarshal([]byte(body), &auth_struct)

	if auth_struct == (Auth{}) || auth_struct.User == (User_struct{}) || auth_struct.Secret == (Secret_struct{}) {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the AuthenticationRequest or it is formed improperly."))
		return
	}

	auth_token := utils.Authenticate(auth_struct.User.Name, auth_struct.Secret.Password)
	fmt.Print("Authentication function passed evaluating now\n")

	if auth_token == "err" {
		writer.WriteHeader(401)
		writer.Write([]byte("The user or password is invalid."))
		return
	}

	writer.WriteHeader(200)
	writer.Write([]byte("\"bearer " + auth_token + "\""))

	fmt.Print("Successfully exiting CreateAuthToken()\n")
}

func GetPackageByName(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("Entering GetPackageByName()\n")

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	name := chi.URLParam(request, "name")

	fmt.Print("Package by name running for name " + name + "\n")

	if given_xAuth == auth_success {

		results := utils.History(name, 0)

		fmt.Print("Results gotten in GetPackageByName()\n")

		if len(results) == 0 {
			writer.WriteHeader(404)
			writer.Write([]byte("No such package."))
			return
		}

		writer.WriteHeader(200)

		for _, i := range results {
			writer.Write([]byte(string(i)))
		}

	} else if given_xAuth == "" || name == "" || given_xAuth != auth_success {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.Write([]byte("{\n  \"code\": 0,\n  \"message\": \"Other Error\"\n}"))
	}

	fmt.Print("Correctly exitting GetPackageByName()\n")
}

func DeletePackageByName(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	name := chi.URLParam(request, "name")

	fmt.Print("Entering DeletePackageByName() for name: " + name + "\n")

	if given_xAuth == auth_success {

		results := utils.History(name, 1)
		fmt.Print("Successfully retreived results in DeletePackageByName()\n")

		if len(results) == 0 {
			writer.WriteHeader(404)
			writer.Write([]byte("Package does not exist."))
			return
		}

		writer.WriteHeader(200)
		writer.Write([]byte("Package is deleted."))

	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

	fmt.Print("Correctly exiting DeletePackageByName() for name: " + name + "\n")
}
