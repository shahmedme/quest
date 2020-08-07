package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"./middleware"
	"./sessions"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", register)
	r.HandleFunc("/secret", middleware.AuthRequired(secretFunc))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":7000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/components/navbar.html", "templates/components/sidebar.html", "templates/partials/meta.html"))
	tmpl.Execute(w, "context")
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func secretFunc(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "secret.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/components/navbar.html", "templates/partials/meta.html"))
		tmpl.Execute(w, "context")
		return
	} else if r.Method == http.MethodPost {
		db, _ := sql.Open("mysql", "root:@(127.0.0.1:3306)/quest?parseTime=true")

		formEmail := r.FormValue("email")
		formPassword := r.FormValue("password")

		var email string
		var password string
		var id int
		var name string

		query := "SELECT * FROM users WHERE email=? and password=?"

		err := db.QueryRow(query, formEmail, formPassword).Scan(&id, &name, &email, &password)

		if err == nil {
			session, _ := sessions.Store.Get(r, "cookie-name")

			session.Values["authenticated"] = true
			session.Save(r, w)

			http.Redirect(w, r, "/secret", http.StatusSeeOther)
		} else {
			log.Fatal(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
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
