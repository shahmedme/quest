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
	templates = template.Must(templates.ParseGlob("templates/partials/*.html"))
	templates = template.Must(templates.ParseGlob("templates/components/*.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/register", register)
	r.HandleFunc("/profile", middleware.AuthRequired(profileFunc))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":7000", r)
}

type User struct {
	Name string
}

func home(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "cookie-name")
	logged := session.Values["authenticated"]

	context := struct {
		Logged bool
		User   *User
	}{Logged: logged.(bool), User: &User{Name: "DIU"}}

	templates.ExecuteTemplate(w, "index.html", context)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func profileFunc(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "profile.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := sessions.Store.Get(r, "cookie-name")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			tmpl := template.Must(template.ParseFiles("templates/login.html", "templates/components/navbar.html", "templates/partials/meta.html"))
			tmpl.Execute(w, "context")
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
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

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			log.Fatal(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := sessions.Store.Get(r, "cookie-name")

		if auth, ok := session.Values["authenticated"].(bool); auth || ok {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		} else {
			tmpl := template.Must(template.ParseFiles("templates/register.html", "templates/components/navbar.html", "templates/partials/meta.html"))
			tmpl.Execute(w, "context")
			return
		}
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
