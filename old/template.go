package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
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
	date      string
	firstname string
	lastname  string
	params    map[string]string
}
type wrapper struct {
	handler http.Handler
}

func (h *wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config := readConfig()
	DB := connectToDatabase(config)
	defer DB.Close()
	context.Set(r, "DB", DB)
	h.handler.ServeHTTP(w, r)
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
	dsn := c.Database["user"] + ":" + c.Database["password"] + "@/" + c.Database["name"]
	con, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Database could not open: %v\n", err)
		os.Exit(1)
	}
	return con
}
func WithDB(h http.Handler) http.Handler {
	return &wrapper{handler: h}
}

type ok interface {
	OK() error
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
	fmt.Println(db)
}
func createLog(w http.ResponseWriter, r *http.Request) {}
func updateLog(w http.ResponseWriter, r *http.Request) {}
func deleteLog(w http.ResponseWriter, r *http.Request) {}

// MAIN

func main() {
	http.Handle("/log", WithDB(http.HandlerFunc(readLog)))
	http.ListenAndServe(":9080", nil)
}
