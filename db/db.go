package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/BenJuan26/samething/config"
	"github.com/BenJuan26/samething/game"
	_ "github.com/lib/pq"
)

var db *sql.DB

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	connStr := fmt.Sprintf("user=%s host=/var/run/postgresql dbname=%s sslmode=disable", config.GetUser(), config.GetDBName())
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Couldn't connect to db: " + err.Error())
	}
}

func GetGameState(id string) (game.State, error) {
	var s game.State
	row := db.QueryRow("SELECT * FROM game WHERE id = '" + id + "'")
	err := row.Scan(&s.ID, &s.State, &s.Player1.Name, &s.Player2.Name, &s.Player1.Word, &s.Player2.Word, &s.Player1.Waiting, &s.Player2.Waiting)
	if err != nil {
		return s, err
	}

	return s, nil
}

func NewGameState() (string, error) {
	valid := false
	var id string
	for !valid {
		id = generateGameID()
		row := db.QueryRow("SELECT COUNT(*) FROM game WHERE id = $1", id)
		var count int64
		err := row.Scan(&count)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if count == 0 {
			valid = true
		}
	}
	_, err := db.Exec("INSERT INTO game VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, 0, "", "", "", "", true, true)
	fmt.Println("New game, id "+id)
	return id, err
}

func UpdateGameState(s game.State) error {
	result, err := db.Exec("UPDATE game SET (state, name1, name2, word1, word2, waiting1, waiting2) = ($1, $2, $3, $4, $5, $6, $7) WHERE id = $8", s.State, s.Player1.Name, s.Player2.Name, s.Player1.Word, s.Player2.Word, s.Player1.Waiting, s.Player2.Waiting, s.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("Incorrect number of rows affected on update: %d", rows)
	}

	return nil
}

func generateGameID() string {
    id := make([]byte, 4)
    for i := range id {
        id[i] = letters[rand.Intn(len(letters))]
    }
    return string(id)
}
