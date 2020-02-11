package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Application struct {
	Client_id    string `json:"client_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Secret       string `json:"secret"`
	Callback_url string `json:"callback_url"`
}

func CreateApplicationsTable() (sql.Result, error) {
	return db.Exec(`CREATE TABLE IF NOT EXISTS applications (
		client_id TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		secret TEXT NOT NULL,
		callback_url TEXT NOT NULL,
		PRIMARY KEY (client_id))`)
}

func AddDummyApplication() {
	dummyApp := Application{Client_id: "thisisaclientid",
		Name:         "Dummy Application",
		Description:  "This is a dummy application",
		Secret:       "hellotherethisisasecret",
		Callback_url: "http://localhost:3000/callback",
	}
	db.Exec("DELETE FROM applications WHERE client_id=$1", dummyApp.Client_id)
	insertApplication(dummyApp)
}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	if !isDigit(r) {
		http.Error(w, "You are not digIT", http.StatusUnauthorized)
		return
	}

	resp, _ := json.Marshal(getAllApplications())
	fmt.Fprint(w, string(resp))
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getAllApplications() []Application {
	rows, err := db.Query("SELECT * FROM applications")
	if err != nil {
		return []Application{}
	}

	var applications []Application
	for rows.Next() {
		var app Application
		err = rows.Scan(
			&app.Client_id,
			&app.Name,
			&app.Description,
			&app.Secret,
			&app.Callback_url)

		if err == nil {
			applications = append(applications, app)
		}
	}

	return applications
}

func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	if !isDigit(r) {
		http.Error(w, "You are not digIT", http.StatusUnauthorized)
		return
	}

	client_id := r.URL.Query().Get("client_id")
	res, _ := db.Exec("DELETE FROM applications WHERE client_id=$1", client_id)
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		http.Error(w, "No applications deleted", http.StatusNotFound)
		return
	}
}
func isDigit(r *http.Request) bool {
	cid := r.URL.Query().Get("cid")
	password := r.URL.Query().Get("password")
	user, err := login_ldap(cid, password)
	if err != nil || !contains(user.Groups, "digit") {
		return false
	}
	return true
}

func AddApplication(w http.ResponseWriter, r *http.Request) {
	if !isDigit(r) {
		http.Error(w, "You are not digIT", http.StatusUnauthorized)
		return
	}

	var app Application

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &app)
	if err != nil {
		http.Error(w, "Could not parse body", http.StatusBadRequest)
		return
	}

	client_id := make([]byte, 32)
	rand.Read(client_id)
	app.Client_id = base64.RawStdEncoding.EncodeToString(client_id)

	secret := make([]byte, 33)
	rand.Read(secret)
	app.Secret = base64.RawStdEncoding.EncodeToString(secret)

	if _, err := insertApplication(app); err != nil {
		http.Error(w, fmt.Sprintf("Unable to insert application: %s", err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, app.toJson())
}

func (app *Application) toJson() string {
	data, _ := json.Marshal(app)
	return string(data)
}

func insertApplication(app Application) (sql.Result, error) {
	return db.Exec("INSERT INTO applications VALUES($1, $2, $3, $4, $5)",
		app.Client_id,
		app.Name,
		app.Description,
		app.Secret,
		app.Callback_url,
	)
}

func CheckClientId(w http.ResponseWriter, r *http.Request) {
	client_id := r.URL.Query().Get("client_id")
	if client_id == "" {
		log.Println("No client id was proviced")
		http.Error(w, "No client id provided", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, _ := json.Marshal(response)
	fmt.Fprint(w, string(resp))
}

func getApplication(client_id string) (Application, error) {
	var app Application
	row := db.QueryRow("SELECT * FROM applications WHERE client_id=$1", client_id)

	err := row.Scan(&app.Client_id,
		&app.Name,
		&app.Description,
		&app.Secret,
		&app.Callback_url,
	)

	return app, err
}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	client_id := r.URL.Query().Get("client_id")

	if client_id == "" {
		log.Println("No client id provided when authenticating")
		http.Error(w, "No client id specified", http.StatusBadRequest)
		return
	}

	var loginParameters struct {
		CID      string
		Password string
	}

	if err := json.Unmarshal(body, &loginParameters); err != nil {
		log.Println("Unable to parse json body")
		log.Println(err)
		http.Error(w, "Could not parse input", http.StatusBadRequest)
		return
	}

	user, err := login_ldap(loginParameters.CID, loginParameters.Password)
	if err != nil {
		log.Printf("Could not log in user: %s\n", loginParameters.CID)
		http.Error(w, "Password and username did not match", http.StatusUnauthorized)
		return
	}

	var app Application
	app, err = getApplication(client_id)
	if err != nil {
		log.Printf("Could not fetch application with id: %s", client_id)
		log.Println(err)
		http.Error(w, "Could not find specified application", http.StatusBadRequest)
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

	resp, _ := json.Marshal(response)
	fmt.Fprint(w, string(resp))
}
