package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("llave"))
}

func principal(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sesion")
	authenticated := session.Values["in"]
	if authenticated == true {
		fmt.Fprintln(w, "Logged!")
	} else {
		fmt.Fprintln(w, "Not Logged!")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sesion")
	session.Values["in"] = true
	session.Values["password"] = 12345
	session.Save(r, w)
	fmt.Fprintln(w, "Bien!")
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sesion")
	session.Values["in"] = false
	session.Values["password"] = 0
	session.Save(r, w)
	fmt.Fprintln(w, "No te vayas!!!! Bueno ni modo")
}

func main() {
	http.HandleFunc("/principal", principal)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe("localhost:8080", nil)

}
