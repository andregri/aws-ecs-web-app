package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "andrea"
	password = "password"
	hostname = "127.0.0.1:3306"
	dbname   = "ecommerce"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createDatabase(db *sql.DB) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	return err
}

func getUserAddress(db *sql.DB, username string) (string, error) {
	query := `select address from user where name = ?`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return "", err
	}
	defer stmt.Close()
	var address string
	row := stmt.QueryRowContext(ctx, username)
	if err := row.Scan(&address); err != nil {
		return "", err
	}
	return address, nil
}

func main() {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	if err := createDatabase(db); err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	http.HandleFunc("/user", userHandler)
	http.ListenAndServe(":9000", nil)
}
