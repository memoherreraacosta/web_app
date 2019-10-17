package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	port      string = "localhost:8080"
	arreglito        = [5]int{1, 5, 25, 125, 5000}
)

type FruitBasket struct {
	Name    string
	Fruit   []string
	Id      int64  `json:"ref"`
	private string // An unexported field is not encoded.
	Created time.Time
}

type Foo struct {
	Bar string
}

func print_array(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "arreglito: %d\n", arreglito[0])
}

func show_html(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "archivo.html")
}

func print_json(w http.ResponseWriter, r *http.Request) {

	basket := FruitBasket{
		Name:    "Standard",
		Fruit:   []string{"Apple", "Banana", "Orange"},
		Id:      999,
		private: "Second-rate",
		Created: time.Now(),
	}

	var jsonData []byte
	jsonData, _ = json.Marshal(basket)

	fmt.Fprint(w, string(jsonData))

}

func download_txt(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "html_textito.txt")
}

func leer_mensaje(w http.ResponseWriter, r *http.Request) {
	foo1 := new(Foo) // or &Foo{}
	json_read := getJson(port+"/json", foo1)
	fmt.Fprintf(w, "json obtenido: %s\n", json_read)

}

func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {

	http.HandleFunc("/arreglito", print_array)
	http.HandleFunc("/pagina", show_html)
	http.HandleFunc("/json", print_json)
	http.HandleFunc("/textito", download_txt)
	http.HandleFunc("/json_leido", leer_mensaje)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		return
	}
}
