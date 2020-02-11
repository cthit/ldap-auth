package main

import (
	"log"
	"os"

	"chalmers.it/ldap-auth/internal/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOrigins:     []string{"https://ldap-auth.chalmers.it", "http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true}))

	router.POST("/api/authenticate", app.HandleAuthenticate)
	router.POST("/api/application", app.AddApplication)
	router.DELETE("/api/application", app.DeleteApplication)
	router.GET("/api/application", app.CheckClientId)
	router.GET("/api/applications", app.GetApplications)

	router.Run(":5011")
}
