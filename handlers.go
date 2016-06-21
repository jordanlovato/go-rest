package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
)

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
	defer stmt.Close()

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
	defer stmt.Close()

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
	defer stmt.Close()

	_, err = stmt.Exec(log_id)
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

// Collections
func readLogs(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*sql.DB)
	var logs []*Log

	rows, err := db.Query("SELECT `firstname`, `lastname`, `date`, `type` FROM `logs` LIMIT 25")
	if err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		log := new(Log)
		if err := rows.Scan(&log.Firstname, &log.Lastname, &log.Date, &log.Type); err != nil {
			respond(w, r, http.StatusBadRequest, err)
			return
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		respond(w, r, http.StatusBadRequest, err)
		return
	}

	respond(w, r, http.StatusOK, &logs)
}
