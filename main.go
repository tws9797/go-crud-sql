package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Album struct {
	ID     int64   `json:"id,omitempty"`
	Title  string  `json:"title,omitempty"`
	Artist string  `json:"artist,omitempty"`
	Price  float32 `json:"price,omitempty"`
}

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var albums []Album
	var err error

	db := dbConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

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

func Read(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var alb Album
	var err error

	db := dbConn()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if err = json.NewDecoder(r.Body).Decode(&alb); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		log.Fatal(err)
	}

	err = db.QueryRow("SELECT * FROM album WHERE id=?", alb.ID).Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)

	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(w).Encode(alb); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		var alb Album
		var id int64
		var err error

		if err = json.NewDecoder(r.Body).Decode(&alb); err != nil {
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			log.Fatal(err)
		}

		db := dbConn()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		dbRes, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		id, err = dbRes.LastInsertId()
		if err != nil {
			log.Fatalf("Unable to get the last inserted ID. %v", err)
		}

		res := response{
			ID:      id,
			Message: "User created successfully",
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Error encoding response object", http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
}

func Update(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		var alb Album
		var err error

		if err = json.NewDecoder(r.Body).Decode(&alb); err != nil {
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			log.Fatal(err)
		}

		db := dbConn()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		_, err = db.Exec("UPDATE album SET title=?,artist=?,price=? WHERE id=?", alb.Title, alb.Artist, alb.Price, alb.ID)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		res := response {
			ID: alb.ID,
			Message: "User updated successfully",
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Error encoding response object", http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		var alb Album
		var err error

		if err = json.NewDecoder(r.Body).Decode(&alb); err != nil {
			http.Error(w, "Error decoding request body", http.StatusInternalServerError)
			log.Fatal(err)
		}

		db := dbConn()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		_, err = db.Exec("DELETE FROM album WHERE id=?", alb.ID)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		res := response {
			ID: alb.ID,
			Message: "User deleted successfully",
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Error encoding response object", http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
}

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/read", Read)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
