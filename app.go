package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/BenJuan26/samething/config"
	"github.com/BenJuan26/samething/db"
	"github.com/BenJuan26/samething/game"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// The idea here is to associate each websocket with a player and a game.
// That will allow the player to reconnect if they've closed their browser
// window or lost their connection. Multiple websockets for the same user
// is allowed in case an old websocket is still hanging around. This is the
// most simple way to do this without some kind of login system.
type Client struct {
	GameID     string
	PlayerName string
	Conn       *websocket.Conn
}

type clientMessage struct {
	Word string `json:"word"`
}

var clients = make([]Client, 0)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var statesChan = make(chan game.State)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", serveIndex)
	router.HandleFunc("/game", newGame).Methods("POST")
	router.HandleFunc("/game/{id}", serveGame)
	router.HandleFunc("/game/{id}/check", checkGame)
	router.HandleFunc("/ws", handleWebsocket)

	go notify()

	http.ListenAndServe(":8080", router)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.ParseFiles("templates/index.html"))

	data := make(map[string]interface{})
	data["base_url"] = config.GetBaseURL()
	data["schema"] = config.GetSchema()
	data["title"] = config.GetTitle()
	indexPage.Execute(w, data)
}

func notify() {
	for {
		state := <-statesChan
		deleted := make([]int, 0)
		for i, client := range clients {
			if client.GameID != state.ID {
				continue
			}
			privState := state

			// Hide the other player's word if the player hasn't submitted theirs yet
			if privState.State == game.WAITING {
				if privState.Player1.Name != client.PlayerName {
					privState.Player1 = game.Player{}
				}
				if privState.Player2.Name != client.PlayerName {
					privState.Player2 = game.Player{}
				}
			}

			err := client.Conn.WriteJSON(privState)
			if err != nil {
				fmt.Println(err)
				deleted = append(deleted, i)
			}
		}

		// Remove any clients we can't connect to
		for numDeleted, i := range deleted {
			j := i - numDeleted
			clients = append(clients[:j], clients[j+1:]...)
		}
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	gameID := query["game"][0]
	name := query["name"][0]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients = append(clients, Client{
		GameID:     gameID,
		PlayerName: name,
		Conn:       conn,
	})

	initialState, err := db.GetGameState(gameID)
	if err != nil {
		fmt.Println(err)
		return
	}
	statesChan <- initialState

	// Listen for new websocket messages forever (or until an error occurs)
	for {
		var msg clientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Addr: %s; Game: %s; Player: %s; Word: %s;\n", conn.RemoteAddr(), gameID, name, msg.Word)
		state, err := db.GetGameState(gameID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonBytes, _ := json.Marshal(state)
		fmt.Println(string(jsonBytes))

		// Set the current player to the first non-empty player.
		// More than two players is an error.
		if state.Player1.Name == "" {
			state.Player1.Name = name
		} else if state.Player2.Name == "" && state.Player1.Name != name {
			state.Player2.Name = name
		} else if state.Player1.Name != name && state.Player2.Name != name {
			fmt.Println("Already 2 players!")
			return
		}

		if state.Player1.Name == name {
			state.Player1.Word = msg.Word
			state.Player1.Waiting = false
		} else if state.Player2.Name == name {
			state.Player2.Word = msg.Word
			state.Player2.Waiting = false
		} else {
			fmt.Println("Unrecognized name " + name)
		}

		state.State = game.WAITING
		if !state.Player1.Waiting && !state.Player2.Waiting {
			fmt.Println("Received both answers")
			if strings.ToLower(state.Player1.Word) == strings.ToLower(state.Player2.Word) {
				fmt.Println("Match!")
				state.State = game.MATCHED
			} else {
				fmt.Println("No match")
				state.State = game.READY
			}
			state.Player1.Waiting = true
			state.Player2.Waiting = true
		}

		err = db.UpdateGameState(state)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Tell the notifier to update the clients about the new game state
		statesChan <- state
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

func checkGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]
	if !db.GameExists(gameID) {
		w.WriteHeader(404)
		return
	}

	fmt.Fprintf(w, `{"redirect_url":"https://benjuan26.com/samething/game/%s"}`, gameID)
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gamePage := template.Must(template.ParseFiles("templates/game.html"))

	data := make(map[string]interface{})
	data["game"] = vars["id"]
	data["base_url"] = config.GetBaseURL()
	gamePage.Execute(w, data)
}
