package main

import (
	"fmt"
	"net/http"

	"github.com/BenJuan26/samething/db"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func (r *http.Request) bool {
        return true;
    },
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}))
	router.HandleFunc("/game", newGame).Methods("POST")
	router.HandleFunc("/game/{id}", serveGame)
	router.HandleFunc("/ws", handleWebsocket)

	http.ListenAndServe(":8080", router)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	game := query["game"][0]
	name := query["name"][0]
	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Println("Game: "+game+"Name: "+name)
	defer conn.Close()
	for {
		_, p, _ := conn.ReadMessage()
		fmt.Println(p)
	}
}

func newGame(w http.ResponseWriter, r *http.Request) {
	id, err := db.NewGameState()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Couldn't create game, try again later")
		return
	}
	fmt.Fprintf(w, `{"redirect_url": "https://benjuan26.com/samething/game/%s"}`, id)
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, `<html><head></head><body><script type="text/javascript">var socket = new WebSocket("wss://benjuan26.com/samething/ws?game=%s&name=%s");socket.send("test");</script></body></html>`, vars["id"], "ben")
}

