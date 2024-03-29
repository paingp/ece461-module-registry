package main

import (
	"log"
	"net/http"
	"os"
	"tomr/frontend"
	"tomr/src/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	// router.Get("/", frontend.RenderHome)
	router.Get("/", frontend.RenderHome)

	router.Get("/UI/packages", frontend.RenderPackages)
	router.Post("/UI/packages", frontend.HandlePackages)

	router.Get("/UI/reset", frontend.RenderReset)
	router.Post("/UI/reset", frontend.HandleReset)

	router.Get("/UI/GETpackage", frontend.RenderGETPackage)
	router.Post("/UI/GETpackage", frontend.HandleGETPackage)

	router.Get("/UI/PUTpackage", frontend.RenderPUTPackage)
	router.Post("/UI/PUTpackage", frontend.HandlePUTPackage)

	router.Get("/UI/DELETEpackage", frontend.RenderDELETEPackage)
	router.Post("/UI/DELETEpackage", frontend.HandleDELETEPackage)

	router.Get("/UI/POSTpackage", frontend.RenderPackage)
	router.Post("/UI/POSTpackage", frontend.HandlePackage)

	router.Get("/UI/Rate", frontend.RenderRate)
	router.Post("/UI/Rate", frontend.HandleRate)

	router.Get("/UI/authenticate", frontend.RenderAuthenticatePackage)
	router.Post("/UI/authenticate", frontend.HandleAuthenticatePackage)

	router.Get("/UI/Regex", frontend.RenderRegex)
	router.Post("/UI/Regex", frontend.HandleRegex)

	/// Additional ////
	router.Get("/UI/Returnhistory", frontend.RenderReturnhistory)
	router.Post("/UI/Returnhistory", frontend.HandleReturnhistory)

	router.Get("/UI/Deletehistory", frontend.RenderDeletehistory)
	router.Post("/UI/Deletehistory", frontend.HandleDeletehistory)

	router.Get("/UI/Listall", frontend.RenderListall)
	router.Post("/UI/Listall", frontend.HandleListall)


	router.Get("/", frontend.RenderHome)
	// router.Get("/", frontend.RenderHome)
	// router.Get("/", frontend.RenderHome)
	// router.Get("/", frontend.RenderHome)

	router.Put("/authenticate", handlers.CreateAuthToken)

	/*
		router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			_, err := writer.Write([]byte("Node Trustworthy Open-source Module Registry -- ECE 461, Team 17"))
			if err != nil {
				log.Println(err)
			}
		}) */

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

	router.Post("/packages", handlers.ListPackages)

	router.Delete("/reset", handlers.ResetRegistry)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		//port = "3000"

	}

	log.Printf("Server started on PORT %s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
