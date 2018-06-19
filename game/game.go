package game

type Player struct {
	Name    string
	Word    string
	Waiting bool
}

type State struct {
	ID      string
	State   int64
	Player1 Player
	Player2 Player
}
