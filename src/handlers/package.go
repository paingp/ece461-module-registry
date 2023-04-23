package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"log"
	"net/http"
	"os/exec"
	"time"

	// "os/exec"
	"path"
	"regexp"
	"tomr/models"
	"tomr/src/db"
	"tomr/src/metrics"
	"tomr/src/utils"

	"github.com/go-chi/chi/v5"
)

const pkgDirPath = "src/metrics/temp" // temp directory to store packages
const auth_success = "ABC"
const bucket_name = "tomr"

type metadata struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	ID      string `json:"ID"`
}

type data struct {
	Conent    string `json:"Content"`
	URL       string `json:"URL"`
	JSProgram string `json:"JSProgram"`
}

type Package struct {
	Metadata metadata `json:"metadata"`
	Data     data     `json:"data"`
}

// const pkgDirPath = "src/metrics/temp" // temp directory to store packages

func CreatePackage(content string, url string, jsprogram string) {
<<<<<<<<< Temporary merge branch 1
=========
	// fmt.Print("Here 222")
>>>>>>>>> Temporary merge branch 2
	packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
	pkgDir := ""
	metadata := models.PackageMetadata{}
	//utils.PrintPackageData(packageData)
	rating := models.PackageRating{}
	var readMe []byte
	// Return Error 400 if both Content and URL are set
	if (packageData.Content != "") && (packageData.URL != "") {
		fmt.Printf("Error 400: Content and URL cannot be both set")
	} else if packageData.Content != "" { // Only Content is set
		// Decode base64 string into zip
		pkgDir = path.Join(pkgDirPath, "package.zip")
		utils.Base64ToZip(packageData.Content, pkgDir)
		err := utils.GetMetadataFromZip(pkgDir, &metadata, &readMe)
		if err != nil {
			log.Fatalf("Failed to get metadata from zip file\n")
		}
		metadata.ID = metadata.Name + "(" + metadata.Version + ")"
		err = metrics.RatePackage(metadata.RepoURL, pkgDir, &rating, metadata.License, &readMe)
	} else { // Only URL is set
		gitUrl := utils.GetGithubUrl(url)
		pkgDir = utils.CloneRepo(gitUrl, pkgDirPath)
		err := metrics.RatePackage(gitUrl, pkgDir, &rating, "", nil)
		if err != nil {
			log.Fatalf("Failed to rate package at URL: %s\n", url)
		}
		// Check if package meets criteria for ingestion
		utils.GetPackageMetadata(pkgDir, &metadata)
		err = utils.ZipDirectory(pkgDir, pkgDir+".zip")
		if err != nil {
			log.Fatal(err)
		}
		pkgDir += ".zip"
	}
	utils.PrintMetadata(metadata)
	utils.PrintRating(rating)
	fmt.Printf(pkgDir)
	// Upload package and store data in system
	pkg := models.PackageObject{Metadata: &metadata, Data: &packageData, Rating: &rating}

	err := db.StorePackage(pkg, pkgDir)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPackageByRegEx(writer http.ResponseWriter, request *http.Request) {

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

		body, _ := ioutil.ReadAll(request.Body)

		var regex regex_body
		json.Unmarshal(body, &regex)

		regex_string = string(body)
	}

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

	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("Status: 400 There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func ResetRegistry(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	if given_xAuth == auth_success {

		cmd := exec.Command("python3", "src/gcp_calc/deleteBucket.py", "tomr")
		_, err := cmd.Output()

		if err != nil {
			println(err.Error())
			return
		}

		writer.WriteHeader(200)
		writer.Write([]byte("Registry is reset."))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

}

func RetrievePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]
	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	id := chi.URLParam(request, "id")

	if given_xAuth == auth_success {

		filepath := "src/handlers/readTo.txt"

		errFile := db.DownloadFile(bucket_name, id, filepath)

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

		attrs, _ := db.GetMetadata(bucket_name, id)
		metadata := attrs.Metadata

		var return_package Package
		return_package.Metadata.Name = rs[1]
		return_package.Metadata.Version = rs[2]
		return_package.Metadata.ID = metadata["ID"]
		return_package.Data.Conent = base64
		return_package.Data.URL = metadata["URL"]
		return_package.Data.JSProgram = metadata["JSProgram"]

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

}

func UpdatePackage(writer http.ResponseWriter, request *http.Request) {

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

		db.DeleteFile(bucket_name, id)

		var recieve_package Package
		body, _ := ioutil.ReadAll(request.Body)
		json.Unmarshal(body, &recieve_package)

		fmt.Print(recieve_package)

		writer.WriteHeader(200)
		writer.Write([]byte("Version is updated."))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func DeletePackage(writer http.ResponseWriter, request *http.Request) {

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

		db.DeleteFile(bucket_name, id)

		writer.WriteHeader(200)
		writer.Write([]byte("Version is deleted."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func ListPackages(writer http.ResponseWriter, request *http.Request) {

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
				writer.WriteHeader(413)
				writer.Write([]byte("Too many packages returned."))
		}
	
	} else if given_xAuth == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.Write([]byte("{\n  \"code\": 0,\n  \"message\": \"Other Error\"\n}"))
	}

}

func RatePackage(writer http.ResponseWriter, request *http.Request) {

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

		// Call go functions here

		writer.WriteHeader(201)
		writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
	} else if given_xAuth == "" || id == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.WriteHeader(401)
	}
}

func CreateAuthToken(writer http.ResponseWriter, request *http.Request) {

	// var auth_string

	// body, _ := ioutil.ReadAll(request.Body)

	// var auth_token auth_body
	// json.Unmarshal(body, &auth_token)

	// fmt.Print(auth_token)

}

func GetPackageByName(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	var given_xAuth string

	if request.Form["X-Authorization"] != nil {
		given_xAuth = request.Form["X-Authorization"][0]

	} else {
		given_xAuth = request.Header["X-Authorization"][0]
	}

	name := chi.URLParam(request, "name")

	if given_xAuth == auth_success {

		results := utils.History(name, 0)

		if len(results) == 0 {
			writer.WriteHeader(404)
			writer.Write([]byte("No such package."))
			return
		}

		writer.WriteHeader(200)

		for _, i := range results {
			writer.Write([]byte(string(i)))
		}

	} else if given_xAuth == "" || name == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.Write([]byte("{\n  \"code\": 0,\n  \"message\": \"Other Error\"\n}"))
	}
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

	if given_xAuth == auth_success {

		results := utils.History(name, 1)

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
}		
