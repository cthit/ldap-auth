package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetApplications(c *gin.Context) {
	if !isDigit(c) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("You are not digIT"))
		return
	}

	c.JSON(http.StatusOK, getAllApplications())
}

func DeleteApplication(c *gin.Context) {
	if !isDigit(c) {
		c.AbortWithError(http.StatusUnauthorized, errors.New("You are not digIT"))
		return
	}

	client_id := c.Request.URL.Query().Get("client_id")
	if !deleteApplication(client_id) {
		c.AbortWithError(http.StatusNotFound, errors.New("No applications deleted"))
		return
	}
}

func AddApplication(c *gin.Context) {
	if !isDigit(c) {
		c.AbortWithError(http.StatusNotFound, errors.New("You are not digIT"))
		return
	}

	var app Application

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &app)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse body"))
		return
	}

	client_id := make([]byte, 32)
	rand.Read(client_id)
	app.Client_id = base64.RawStdEncoding.EncodeToString(client_id)

	secret := make([]byte, 33)
	rand.Read(secret)
	app.Secret = base64.RawStdEncoding.EncodeToString(secret)

	if _, err := insertApplication(app); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, app)
}

func CheckClientId(c *gin.Context) {
	client_id := c.Request.URL.Query().Get("client_id")
	if client_id == "" {
		log.Println("No client id was proviced")
		c.AbortWithError(http.StatusBadRequest, errors.New("No client id provided"))
		return
	}

	row := db.QueryRow("SELECT client_id, name, description FROM applications WHERE client_id=$1", client_id)

	var response struct {
		Client_id   string `json:"client_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := row.Scan(
		&response.Client_id,
		&response.Name,
		&response.Description)
	if err != nil {
		log.Printf("Could not authenticate user for client: %s\n", client_id)
		log.Println(err)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func HandleAuthenticate(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	client_id := c.Request.URL.Query().Get("client_id")

	if client_id == "" {
		log.Println("No client id provided when authenticating")
		c.AbortWithError(http.StatusBadRequest, errors.New("No client id specified"))
		return
	}

	var loginParameters struct {
		CID      string
		Password string
	}

	if err := json.Unmarshal(body, &loginParameters); err != nil {
		log.Println("Unable to parse json body")
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse input"))
		return
	}

	user, err := login_ldap(loginParameters.CID, loginParameters.Password)
	if err != nil {
		log.Printf("Could not log in user: %s\n", loginParameters.CID)
		c.AbortWithError(http.StatusUnauthorized, errors.New("Password and username did not match"))
		return
	}

	var app Application
	app, err = getApplication(client_id)
	if err != nil {
		log.Printf("Could not fetch application with id: %s", client_id)
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("Could not find specified application"))
		return
	}

	user.ExpiresAt = time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	var response struct {
		Callback_url string `json:"callback_url"`
		Token        string `json:"token"`
	}

	response.Callback_url = app.Callback_url
	response.Token, _ = token.SignedString([]byte(app.Secret))

	c.JSON(http.StatusOK, response)
}

func isDigit(c *gin.Context) bool {
	headerAuth := c.GetHeader("Authorization")
	if !strings.HasPrefix(headerAuth, "Basic") {
		return false
	}

	decoded, err := base64.RawStdEncoding.DecodeString(strings.Split(headerAuth, " ")[1])
	if err != nil {
		return false
	}

	cidAndPass := strings.Split(string(decoded), ":")

	user, err := login_ldap(cidAndPass[0], cidAndPass[1])
	if err != nil || !contains(user.Groups, "digit") {
		return false
	}
	return true
}
