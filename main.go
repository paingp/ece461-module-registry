package main

import (
	"log"
	"tomr/src/handlers"
	"tomr/src/utils"
)

// "fmt"

func main() {
	//url := "https://www.npmjs.com/package/du"

	content, err := utils.ZipToBase64("src/metrics/temp/node-du.zip")
	if err != nil {
		log.Fatal(err)
	}

	handlers.CreatePackage(content, "", "console.log('Hello World')")

	/*
		bucketName := "tomr"
		err := db.CreateBucket(bucketName)
		if err != nil {
			panic(err)
		}
	*/
	/*
		pkgDir := "src/metrics/temp/node-du.zip"
		objectName := "node-du(^1.9.4)"
		err := db.UploadPackage(pkgDir, "tomr", objectName)
		if err != nil {
			panic(err)
		}
	*/
	/*
		router := chi.NewRouter()

		router.Put("/authenticate", handlers.CreateAuthToken)

		router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			_, err := writer.Write([]byte("ece461g17-module-registry"))
			if err != nil {
				log.Println(err)
			}
		})

		router.Route("/package", func(r chi.Router) {
			// r.Post("/", handlers.CreatePackage)
			r.Get("/{id}", handlers.RetrievePackage)
			r.Put("/{id}", handlers.UpdatePackage)
			r.Delete("/{id}", handlers.DeletePackage)
			r.Get("/{id}/rate", handlers.RatePackage)
		})

		router.Get("/package/byName/{name}", handlers.GetPackageByName)
		router.Delete("/package/byName/{name}", handlers.DeletePackageByName)

		router.Post("/package/byRegEx", handlers.GetPackageByRegEx)

		router.Post("/packages", handlers.ListPackages)

		router.Delete("/reset", handlers.ResetRegistry)

		err := http.ListenAndServe(":3000", router)
		if err != nil {
			log.Println(err)
		}
	*/
}
