package game

const (
    STARTED = 0
    WAITING = 1
    READY   = 2
    MATCHED = 3
)

type Player struct {
	Name    string `json:"name"`
	Word    string `json:"word"`
	Waiting bool   `json:"waiting"`
}

type State struct {
	ID      string `json:"id"`
	State   int64  `json:"state"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
}
