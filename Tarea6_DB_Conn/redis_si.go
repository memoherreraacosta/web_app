package main

import (
	"fmt"
	"log"
	"net/http"

  redisStore "gopkg.in/boj/redistore.v1"
  
)

var store *redisStore.RediStore
var problema error

func init() {
	store, problema = redisStore.NewRediStore(10, "tcp", "redis-12257.c57.us-east-1-4.ec2.cloud.redislabs.com:12257", "00T596eexiAEkuKOYT9QNJsfsFaMTfOm", []byte("secret-key"))
	if problema != nil {
		log.Fatal("Redis Error : ", problema)
	}
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
	session.Save(r, w)
	fmt.Fprintln(w, "Bien!")
}
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sesion")
	session.Values["in"] = false
	session.Save(r, w)
	fmt.Fprintln(w, "No te vayas!!!! Bueno ni modo")
}
func main() {
	http.HandleFunc("/principal", principal)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe("localhost:8080", nil)
	defer store.Close()

}
