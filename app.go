package main

import (
	"fmt"

	"github.com/BenJuan26/samething/db"
	"github.com/BenJuan26/samething/game"
)

func main() {
	state := game.State{
		ID:    "ABCD",
		Player1:game.Player{Word: "banana"},
	}
	err := db.UpdateGameState(state)
	if err != nil {
		panic(err)
	}
	s, err := db.GetGameState("ABCD")
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Player1.Word)
}
