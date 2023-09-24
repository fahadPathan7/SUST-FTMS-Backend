package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ftms/models"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// connecting to mysql database
func CreateDbConnection() {
	var err error
	db, err = sql.Open("mysql", "fahadftms:fahadftms@tcp(localhost:3306)/ftms") // port 3306 is the default port for mysql in xampp
	// here ftms is the database name

	if err != nil {
		fmt.Println("Error connecting databse!")
		panic(err.Error())
	}

	// defer db.Close()
	fmt.Println("Successfully connected to mysql database")
}





// insert dept info into database
func insertNewDept(dept models.Dept) {
	// dept.DeptCode is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbldept VALUES (?, ?, ?)", dept.DeptCode, dept.DeptName, dept.DeptShortName)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// check if dept exists in database
func deptExists(deptCode int) bool {
	var dept models.Dept
	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptShortName)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// controller function to insert new dept
func InsertNewDept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	if !deptExists(dept.DeptCode) {
		insertNewDept(dept)
		json.NewEncoder(w).Encode(dept)
	} else {
		json.NewEncoder(w).Encode("Dept already exists!")
	}
}