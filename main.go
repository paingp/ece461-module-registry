// package main

// import (
// 	"log"
// 	"tomr/src/handlers"
// 	"tomr/src/utils"
// )

// func main() {

// 	//fmt.Println(os.Getenv("GITHUB_TOKEN"))

// 	//url := "https://www.npmjs.com/package/du"
// 	//gitUrl := utils.GetGithubUrl(url)
// 	//utils.CloneRepo(gitUrl, "src/metrics/temp")

// 	/*
// 		pkgDir := "src/metrics/temp/package.zip"
// 		metadata := models.PackageMetadata{}
// 		var readme []byte
// 		utils.GetMetadataFromZip(pkgDir, &metadata, &readme)
// 		fmt.Printf(string(readme))
// 	*/
// 	//pkgDir := ""

// 	// Encode/decode between base64 string and ZIP

// 	content, err := utils.ZipToBase64("src/metrics/temp/node-du.zip")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	/*
// 		err = utils.Base64ToZip(content, "src/metrics/temp/package.zip")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	*/
// 	handlers.CreatePackage(content, "", "console.log('Hello World')")

// 	//metrics.RatePackage(url, pkgDir)

// }




package main

import (
	// "fmt"

	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	auth_success := "ABC"

	router := chi.NewRouter()

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("ece461g17-module-registry"))
		if err != nil {
			log.Println(err)
		}
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Println(err)
	}

	// /package/byRegEx:
	router.Post("/package/byRegEx", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		given_xAuth := request.Form["X-Authorization"][0]
		regex_string := request.Form["Regex String"][0]

		if given_xAuth == auth_success && regex_string != "" {

			// Calls Go Function here

			writer.WriteHeader(200)
			writer.Write([]byte("Status: 200 Return a list of packages."))

		} else if regex_string == "" || given_xAuth == "" {
			writer.WriteHeader(400)
			writer.Write([]byte("Status: 400 There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		}
	})

	// /reset:
	router.Delete("/reset", func(writer http.ResponseWriter, request *http.Request) {

		request.ParseForm()
		given_xAuth := request.Form["X-Authorization"][0]

		fmt.Print((given_xAuth))

		if given_xAuth == auth_success {

			// Call Go Function here 
			fmt.Print("Aditya Srikanth")

			writer.WriteHeader(200)
			writer.Write([]byte("Registry is reset."))
		} else if given_xAuth == ""{
			writer.WriteHeader(400)
			writer.Write([]byte("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		} else {
			writer.WriteHeader(401)
			writer.Write([]byte("You do not have permission to reset the registry."))
		}

	})

	// get /package/{id}
	router.Get("/package/{id}", func(writer http.ResponseWriter, request *http.Request) {

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
	})

	// put /package/{id}
	router.Put("/package/{id}", func(writer http.ResponseWriter, request *http.Request) {

		request.ParseForm()

		given_xAuth := request.Form["X-Authorization"][0]
		id := chi.URLParam(request, "username")

		if given_xAuth == auth_success {

			fmt.Print(id)

			// Call go functions here 
			
			writer.WriteHeader(200)
			writer.Write([]byte("Version is updated."))
		} else {
			writer.WriteHeader(401)
			writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		}
	})

	// delete /package/{id}
	router.Delete("/package/{id}", func(writer http.ResponseWriter, request *http.Request) {

		request.ParseForm()

		given_xAuth := request.Form["X-Authorization"][0]
		id := chi.URLParam(request, "username")

		if given_xAuth == auth_success {

			fmt.Print(id)

			// Call go functions here 
			
			writer.WriteHeader(200)
			writer.Write([]byte("Version is deleted."))
		} else {
			writer.WriteHeader(401)
			writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		}
	})

	// /package 
	router.Post("/package", func(writer http.ResponseWriter, request *http.Request){
		
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

	})

	// get /package/{id}/rate
	router.Get("/package/{id}/rate", func(writer http.ResponseWriter, request *http.Request) {

		request.ParseForm()

		given_xAuth := request.Form["X-Authorization"][0]
		id := chi.URLParam(request, "username")

		if given_xAuth == auth_success {

			fmt.Print(id)

			// Call go functions here 

			writer.WriteHeader(201)
			writer.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
		} else {
			writer.WriteHeader(401)
			writer.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		}
	})

	// /authenticate
	router.Get("/authenticate", func(writer http.ResponseWriter, request *http.Request) {

	})

	router.Get("/package/byName/{name}" , func(writer http.ResponseWriter, request *http.Request) {
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
	})

	router.Delete("/package/byName/{name}" , func(writer http.ResponseWriter, request *http.Request) {
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
	})

}
