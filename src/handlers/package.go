package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"tomr/models"
	"tomr/src/metrics"
	"tomr/src/utils"

	"github.com/go-chi/chi/v5"
)

const pkgDirPath = "src/metrics/temp" // temp directory to store packages
const auth_success = "ABC"

func CreatePackage(content string, url string, jsprogram string) {
	packageData := models.PackageData{Content: content, URL: url, JSProgram: jsprogram}
	//utils.PrintPackageData(packageData)
	// Return Error 400 if both Content and URL are set
	if (packageData.Content != "") && (packageData.URL != "") {
		fmt.Printf("Error 400: Content and URL cannot be both set")
	} else if packageData.Content != "" { // Only Content is set
		// Decode base64 string into zip
		pkgDir := path.Join(pkgDirPath, "package.zip")
		utils.Base64ToZip(packageData.Content, pkgDir)
		var readMe []byte
		metadata, err := utils.GetMetadataFromZip(pkgDir, &readMe)
		if err != nil {
			log.Fatalf("Failed to get metadata from zip file\n")
		}
		utils.PrintMetadata(metadata)
		metrics.RatePackage(metadata.RepoURL, pkgDir, metadata.License, &readMe)
	} else {
		gitUrl := utils.GetGithubUrl(url)
		pkgDir := utils.CloneRepo(gitUrl, pkgDirPath)
		metrics.RatePackage(gitUrl, pkgDir, "", nil)
	}
	// Upload package and store data in system
}

// ///////////////////////CreatePackage///////////////////////////
func GetPackageByRegEx(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]
	regex_string := request.Form["Regex String"][0]

	if given_xAuth == auth_success && regex_string != "" {

		writer.WriteHeader(200)
		writer.Write([]byte("Status: 200 Return a list of packages."))

	} else if regex_string == "" || given_xAuth == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("Status: 400 There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func ResetRegistry(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	given_xAuth := request.Form["X-Authorization"][0]

	fmt.Print((given_xAuth))

	if given_xAuth == auth_success {

		// Call Go Function here
		fmt.Print("Aditya Srikanth")

		writer.WriteHeader(200)
		writer.Write([]byte("Registry is reset."))
	} else if given_xAuth == "" {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("You do not have permission to reset the registry."))
	}
}

func RetrievePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]
	id := chi.URLParam(request, "username")

	if given_xAuth == auth_success {

		fmt.Print(id)

		// Call go functions here

		writer.WriteHeader(200)
		writer.Write([]byte("Return the package. Content is required."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func UpdatePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]
	id := chi.URLParam(request, "id")

	if given_xAuth == auth_success {

		fmt.Print(id)

		// Call go functions here

		writer.WriteHeader(200)
		writer.Write([]byte("Version is updated."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func DeletePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]
	id := chi.URLParam(request, "id")

	if given_xAuth == auth_success {

		fmt.Print(id)

		// Call go functions here

		writer.WriteHeader(200)
		writer.Write([]byte("Version is deleted."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func ListPackages(writer http.ResponseWriter, request *http.Request) {

	fmt.Print("made it to ListPackages")

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]

	if given_xAuth == auth_success {

		// Call go functions here

		writer.WriteHeader(201)
		writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
	} else {
		writer.WriteHeader(400)
		writer.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}

}

func RatePackage(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	given_xAuth := request.Form["X-Authorization"][0]
	id := chi.URLParam(request, "id")

	if given_xAuth == auth_success {

		fmt.Print(id)

		// Call go functions here

		writer.WriteHeader(201)
		writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func CreateAuthToken(writer http.ResponseWriter, request *http.Request) {

}

func GetPackageByName(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	given_xAuth := request.Form["X-Authorization"][0]

	if given_xAuth == auth_success {

		// Call go functions here

		writer.WriteHeader(200)
		writer.Write([]byte("Return the package history."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

func DeletePackageByName(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	given_xAuth := request.Form["X-Authorization"][0]

	if given_xAuth == auth_success {

		// Call go functions here

		writer.WriteHeader(200)
		writer.Write([]byte("Return the package history."))
	} else {
		writer.WriteHeader(401)
		writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	}
}

////////////////////////////////////////////////////

/*
func GetPackageMetadata1(directory string, isZip bool) {
	pkgJsonPath := ""
	var rc io.ReadCloser
	//metadata := models.PackageMetadata{}
	if isZip {
		pkgJsonPath = utils.GetMetadataFromZip(directory)
	} else {
		pkgJsonPath = path.Join(pkgJsonPath, "/package.json")
	}
	rc, err := os.Open(pkgJsonPath)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(rc)
	type metadata struct {
		Name    string `json:"Name"`
		Version string `json:"Version"`
		License string `json:"License"`
	}
	var m metadata
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	rc.Close()
	fmt.Println(m.Name)
	fmt.Println(m.Version)
	fmt.Println(m.License)

		data, err := os.ReadFile(pkgJsonPath)
		if err != nil {
			log.Fatal(err)
		}

		//packageMetadata := models.PackageMetadata{}
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(data), &jsonMap)

		metadata := models.PackageMetadata{Name: jsonMap["name"].(string), Version: jsonMap["version"].(string), ID: ""}
		metadata.ID = metadata.Name + "_" + metadata.Version

	//utils.PrintMetadata(metadata)
}
*/
