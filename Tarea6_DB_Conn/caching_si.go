package main

import (
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var newCache *cache.Cache

func init() {
	newCache = cache.New(1*time.Minute, 10*time.Minute)

}
func darRespuesta(w http.ResponseWriter, r *http.Request) {
	dato, existe := newCache.Get("cancion")

}
func main() {
	http.HandleFunc("/", darRespuesta)
	http.ListenAndServe("localhost:8080", nil)

}
