package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func dbConn() (db *sql.DB) {

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")

	return db
}

func Index(w http.ResponseWriter, r *http.Request) {

	var albums []Album

	db := dbConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	w.Header().Add("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			log.Fatal(err)
		}
		albums = append(albums, alb)
	}

	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {

	var album Album

	db := dbConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}

}

func main() {

	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
