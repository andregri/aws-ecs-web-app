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
	dbname   = "phonebook"
)

var db *sql.DB

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Get from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address, err := getUserAddress(db, name)
		if err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Address = %s\n", address)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		err := createUser(db, name, address)
		if err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
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

func createTable(db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS user(id int primary key auto_increment, name text, address text)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.ExecContext(ctx, query)
	return err
}

func createUser(db *sql.DB, name string, address string) error {
	query := "INSERT INTO user(name, address) VALUES (?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, name, address)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	return nil
}

func getUserAddress(db *sql.DB, name string) (string, error) {
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
	row := stmt.QueryRowContext(ctx, name)
	if err := row.Scan(&address); err != nil {
		return "", err
	}
	return address, nil
}

func main() {
	var err error
	db, err = sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	if err := createDatabase(db); err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}

	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Printf("Error %s when creating table\n", err)
		return
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	http.HandleFunc("/user", userHandler)
	http.ListenAndServe(":9000", nil)
}
