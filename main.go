package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tomr/src/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Put("/authenticate", handlers.CreateAuthToken)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("ece461g17-module-registry"))
		if err != nil {
			log.Println(err)
		}
	})

	router.Route("/package", func(r chi.Router) {
		r.Post("/", handlers.CreatePackage)
		r.Get("/{id}", handlers.RetrievePackage)
		r.Put("/{id}", handlers.UpdatePackage)
		r.Delete("/{id}", handlers.DeletePackage)
		r.Get("/{id}/rate", handlers.RatePackage)
	})

	router.Get("/package/byName/{name}", handlers.GetPackageByName)
	router.Delete("/package/byName/{name}", handlers.DeletePackageByName)

	router.Post("/package/byRegEx", handlers.GetPackageByRegEx)

	router.Post("/`packages`", handlers.ListPackages)

	router.Delete("/reset", handlers.ResetRegistry)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started on PORT %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
