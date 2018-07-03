package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/employees", EmployeeMS)
	r.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":8000", r))
}

//Index function
func Index(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	day := t.Format("2006 January _2 03:04:05PM MST")
	fmt.Fprintf(w, "<h1>Employee service is alive.<br>The current time is: "+day+"</h1>")
}

//EmployeeMS will invoke the employee microservice
func EmployeeMS(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, EmployeeDb(10))
}

//EmployeeDb is a database that holds employees
func EmployeeDb(limit int) string {

	var (
		dbName     = os.Getenv("MYSQL_DATABASENAME") //employees
		dbHost     = os.Getenv("MYSQL_DATABASEHOST") //192.168.64.3
		dbPort     = os.Getenv("MYSQL_PORT")         //31734
		dbUser     = os.Getenv("MYSQL_USER")         //root
		dbPassword = os.Getenv("MYSQL_PASSWORD")     //mysql
	)

	type EmployeeData struct {
		DeptName      string
		EmpNo         int
		FirstName     string
		LastName      string
		Title         string
		TitleFromDate string
		TitleToDate   string
		Gender        string
		BirthDate     string
		HireDate      string
	}
	var ed EmployeeData
	var eds []EmployeeData

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(
		"SELECT d.dept_name, .e.emp_no," +
			"e.first_name, e.last_name, t.title, DATE_FORMAT(t.from_date,'%Y-%M-%D %T'), DATE_FORMAT(t.to_date, '%Y-%M-%D %T'), " +
			"e.gender, DATE_FORMAT(e.birth_date,'%Y-%M-%D %T'), DATE_FORMAT(e.hire_date, '%Y-%M-%D %T') " +
			"FROM employees as e " +
			"INNER JOIN dept_emp as de ON e.emp_no = de.emp_no " +
			"INNER JOIN dept_manager AS dm ON dm.emp_no = e.emp_no " +
			"INNER JOIN departments AS d ON d.dept_no = dm.dept_no " +
			"INNER JOIN titles as t ON e.emp_no = t.emp_no " +
			"ORDER BY e.last_name, e.first_name, e.hire_date, t.from_date DESC")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	for rows.Next() {
		err := rows.Scan(&ed.DeptName, &ed.EmpNo, &ed.FirstName, &ed.LastName, &ed.Title, &ed.TitleFromDate, &ed.TitleToDate, &ed.Gender, &ed.BirthDate, &ed.HireDate)
		checkErr(err)

		eds = append(eds, ed)

	}
	strJSON, err := json.Marshal(eds)
	checkErr(err)
	return string(strJSON)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
