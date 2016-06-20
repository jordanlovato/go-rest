package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/pat"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// TYPES
type ConfigMap struct {
	Database map[string]string
}
type Log struct {
	Date      string
	Firstname string
	Lastname  string
	Type      string
}

func WithDB(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config := readConfig()
		DB := connectToDatabase(config)
		defer DB.Close()
		context.Set(r, "DB", DB)
		fn(w, r)
	}
}
func readConfig() *ConfigMap {
	configFile := "config.json"
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var config ConfigMap
	json.Unmarshal(file, &config)
	return &config
}
func connectToDatabase(c *ConfigMap) *sql.DB {
	dsn := c.Database["username"] + ":" + c.Database["password"] + "@/" + c.Database["name"]
	con, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Database could not open: %v\n", err)
		os.Exit(1)
	}
	return con
}

type ok interface {
	OK() error
}

func (l *Log) OK() error {
	// Basic validation
	if len(l.Firstname) == 0 {
		return ErrRequired("Firstname")
	}

	if len(l.Lastname) == 0 {
		return ErrRequired("Lastname")
	}

	return nil
}

func ErrRequired(field string) error {
	fmt.Println("%s", field)
	return nil
}
func decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	if validatable, ok := v.(ok); ok {
		return validatable.OK()
	}
	return nil
}
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("respond:", err)
	}
}

// HANDLERS
func readLog(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*sql.DB)
	var log Log
	var log_id = r.URL.Query().Get(":log_id")

	err := db.QueryRow("SELECT `firstname`, `lastname`, `date`, `type`  FROM `logs` WHERE `id`=?", log_id).Scan(&log.Firstname, &log.Lastname, &log.Date, &log.Type)
	if err != nil {
		fmt.Println("get rekt:", err)
	}
	respond(w, r, http.StatusOK, &log)
}
func createLog(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*sql.DB)
	var log Log

	if err := decode(r, &log); err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	stmt, err := db.Prepare("INSERT logs SET `firstname`=?,`lastname`=?,`date`=?,`type`=?")
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	_, err = stmt.Exec(log.Firstname, log.Lastname, log.Date, log.Type)
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusOK, &log)
}
func updateLog(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*sql.DB)
	var log Log
	var log_id = r.URL.Query().Get(":log_id")

	if err := decode(r, &log); err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	stmt, err := db.Prepare("UPDATE logs SET `firstname`=?,`lastname`=?,`date`=?,`type`=? WHERE `id`=?")
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	_, err = stmt.Exec(log.Firstname, log.Lastname, log.Date, log.Type, log_id)
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusOK, &log)
}
func deleteLog(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*sql.DB)
	var log_id = r.URL.Query().Get(":log_id")

	stmt, err := db.Prepare("DELETE FROM `logs` WHERE `id`=?")
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	_, err = stmt.Exec(log_id)
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

// MAIN

func main() {
	r := pat.New()
	r.Get("/log/{log_id}", WithDB(http.HandlerFunc(readLog)))
	r.Put("/log", WithDB(http.HandlerFunc(createLog)))
	r.Patch("/log/{log_id}", WithDB(http.HandlerFunc(updateLog)))
	r.Delete("/log/{log_id}", WithDB(http.HandlerFunc(deleteLog)))

	http.ListenAndServe(":9080", context.ClearHandler(r))
}
