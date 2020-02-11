package main

import (
	"log"
	"net/http"
	"os"

	"chalmers.it/ldap-auth/internal/app"
	"github.com/rs/cors"
)

func main() {
	if os.Getenv("MOCK_MODE") == "true" {
		app.SetupMock()
		return
	}

	app.SetupDB()
	defer app.CloseDB()

	app.CreateApplicationsTable()
	if os.Getenv("ADD_DUMMY_APP") == "true" {
		app.AddDummyApplication()
	}

	log.Println("Starting")
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"https://ldap-auth.chalmers.it", "http://localhost:3000"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true})

	mux.HandleFunc("/api/authenticate", app.HandleAuthenticate)
	mux.HandleFunc("/api/application/add", app.AddApplication)
	mux.HandleFunc("/api/application/delete", app.DeleteApplication)
	mux.HandleFunc("/api/application", app.CheckClientId)
	mux.HandleFunc("/api/applications", app.GetApplications)

	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":5011", handler))
}
