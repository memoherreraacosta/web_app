package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Libro struct {
	Id     string `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

type Libros []Libro

var libros []Libro

func init() {
	libros = Libros{
		Libro{Id: "1", Titulo: "Las ciudades invisibles", Autor: "Italo Calvino"},
		Libro{Id: "2", Titulo: "Jerusalem", Autor: "Mia Couto"},
		Libro{Id: "3", Titulo: "Las Cosas", Autor: "Goerges Perec"},
	}
}

func getLibro(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//fmt.Printf("%s\n", vars["id"])

	id := vars["id"]
	for _, l := range libros {
		if l.Id == id {
			if e := json.NewEncoder(w).Encode(l); e != nil {
				fmt.Printf("Error :: %e", e)
			}
		}
		//	fmt.Println(r.Method)
	}
}

func getLibros(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// fmt.Printf("%s\n", r.Method)

	id := vars["id"]
	idInt, _ := strconv.Atoi(id)
	// fmt.Printf("%d\n", idInt)
	// fmt.Printf("%d\n", len(libros))

	if e := json.NewEncoder(w).Encode(libros); e != nil &&
		r.Method == "GET" &&
		0 < idInt &&
		idInt < len(libros) {
		fmt.Printf("Error :: %e", e)
	}
}

func anadirLibro(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%s\n", r.Method)
	libro := Libro{}

	json.NewDecoder(r.Body).Decode(&libro)
	if r.Method == "POST" {
		libros = append(libros, libro)
	}
}

func updateLibro(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PUT" {
		fmt.Println("Not a put method")
		return
	}
	libro := Libro{}
	json.NewDecoder(r.Body).Decode(&libro)

	inserted := false
	for index, l := range libros {
		if l.Id == libro.Id {
			libros[index] = libro
			inserted = !inserted
		}
	}
	if !inserted {
		libros = append(libros, libro)
	}
}

func deleteLibro(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		fmt.Println("Not a delete method")
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]

	deleted := false
	for index, l := range libros {
		if l.Id == id {
			libros[index] = libros[len(libros)-1]
			libros = libros[:len(libros)-1]
			deleted = !deleted
		}
	}
	if !deleted {
		fmt.Printf("Book not found by its id : %s\n", id)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/libros/{id}", getLibro)        // GET method
	r.HandleFunc("/todos_libros/{id}", getLibros) // GET method
	r.HandleFunc("/add_libro", anadirLibro).Methods("POST")
	r.HandleFunc("/update_libro", updateLibro).Methods("PUT")
	r.HandleFunc("/delete_libro/{id}", deleteLibro).Methods("DELETE")
	/*
		{"id":"4","titulo":"Visión de los Vencidos","autor":"Miguel León-Portilla"}
	*/
	http.ListenAndServe(":8080", r)
}
