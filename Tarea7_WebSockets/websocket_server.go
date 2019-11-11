package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}

func HandleClients(w http.ResponseWriter, r *http.Request) {

	go broadcastMessages()

	//"Mejora" la conexion
	websocket, _ := upgrader.Upgrade(w, r, nil)

	defer websocket.Close()

	clients[websocket] = true

	for {
		var m Message
		err := websocket.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error leyendo el mensaje")
			delete(clients, websocket)
			break
		}

		fmt.Println(m)
		//Enviar al canal broadcast
		broadcast <- m
	}
}

func principal(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index_w.html")
}

func main() {
	http.HandleFunc("/", principal)
	http.HandleFunc("/echo", HandleClients)
	http.ListenAndServe(":8081", nil)

}

func broadcastMessages() {
	for {
		//Recibir y guardar
		message := <-broadcast
		fmt.Println(message)
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				fmt.Println("Error al enviar al canal")
				client.Close()
				delete(clients, client)
			}
		}
	}
}
