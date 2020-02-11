package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

const mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiJ0ZXN0VXNlciIsIm5pY2siOiJOaWNrbmFtZSIsImdyb3VwcyI6WyJkaWdpdCIsInByaXQiLCJzdHlyaXQiXX0.EwoDK_VMgDhjLTpJTku9KRDZB4-tMwLqaSCgMHzVAkI"

var dummyApp = Application{Client_id: "thisisaclientid",
	Name:         "Dummy Application",
	Description:  "This is a dummy application",
	Secret:       "hellotherethisisasecret",
	Callback_url: "http://localhost:3000/callback",
}

func SetupMock() {
	log.Println("Starting")
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:3011"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true})

	mux.HandleFunc("/api/authenticate", mockAuth)
	mux.HandleFunc("/api/application", mockCheckId)

	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":5011", handler))
}

func mockAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprintf(`{"callback_url":"http://localhost:3000/callback", "token": "%s"}`, mockToken))
}

func mockCheckId(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("client_id") != dummyApp.Client_id {
		http.Error(w, "[Mock mode] Application does not exit", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, `{
		"client_id": "thisisaclientid",
		"name": "Dummy Application",
		"description": "This is a dummy application"
	  }`)
}
