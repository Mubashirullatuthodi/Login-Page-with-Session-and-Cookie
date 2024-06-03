package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/icza/session"
)

var tmpl *template.Template

func init() {

	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}
func clearCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	fmt.Println("URL:= ", r.URL)
	sess := session.Get(r)
	if sess == nil {
		title := map[string]interface{}{
			"head": "LOGIN PAGE",
		}

		tmpl.ExecuteTemplate(w, "index.html", title)
	} else {
		http.Redirect(w, r, "/Home", http.StatusSeeOther)
	}
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	fmt.Println("Method: ", r.Method)
	fmt.Println("URL:= ", r.URL)
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "Mubashir" && password == "12345678" {

		sess := session.NewSessionOptions(&session.SessOptions{
			CAttrs: map[string]interface{}{"username": username, "password": password}, //adding sessions
		})
		session.Add(sess, w)
		http.Redirect(w, r, "/Home", http.StatusSeeOther)

	} else {
		data := map[string]interface{}{
			"error": "Invalid Username and Password",
			"head":  "LOGIN PAGE",
		}
		tmpl.ExecuteTemplate(w, "index.html", data)
	}

}
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	fmt.Println("URL:= ", r.URL)
	sess := session.Get(r)
	if sess != nil {

		username := sess.CAttr("username")

		data := map[string]interface{}{
			"username": username,
			"head":     "HOME PAGE",
		}
		tmpl.ExecuteTemplate(w, "Home.html", data)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	fmt.Println("URL:= ", r.URL)
	sess := session.Get(r)
	if sess != nil {
		session.Remove(sess, w)
		fmt.Println("LOGOUT")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func main() {
	 fmt.Println("STARTING THE SERVER :8080")
	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/Home", HomePageHandler)
	http.HandleFunc("/logout", LogOutHandler)
	http.ListenAndServe(":8080", nil)

}
