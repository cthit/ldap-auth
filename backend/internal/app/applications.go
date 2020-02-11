package app

import (
	"database/sql"
	"encoding/json"
)

type Application struct {
	Client_id    string `json:"client_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Secret       string `json:"secret"`
	Callback_url string `json:"callback_url"`
}

func (app *Application) toJson() string {
	data, _ := json.Marshal(app)
	return string(data)
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

func deleteApplication(client_id string) bool {
	res, _ := db.Exec("DELETE FROM applications WHERE client_id=$1", client_id)
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return false
	}
	return true
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
