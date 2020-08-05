package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":7000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/components/navbar.html", "templates/components/sidebar.html", "templates/partials/meta.html"))
	tmpl.Execute(w, "context")
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/components/navbar.html", "templates/partials/meta.html"))
	tmpl.Execute(w, "context")
	return
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/register.html", "templates/components/navbar.html", "templates/partials/meta.html"))
		tmpl.Execute(w, "context")
		return
	} else if r.Method == http.MethodPost {
		db, _ := sql.Open("mysql", "root:@(127.0.0.1:3306)/quest?parseTime=true")

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		query := `INSERT INTO users(name, email, password) VALUES(?, ?, ?)`
		_, err := db.Exec(query, name, email, password)

		if err == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			fmt.Print(err)
		}
	}
}
