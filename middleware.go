package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
)

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
		return nil
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
		return nil
	}
	return con
}
