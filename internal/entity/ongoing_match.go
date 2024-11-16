package entity

type OngoingMatch struct {
	MatchGame          Match
	ScoreCurrentPlayer Score
	ScoreAnotherPlayer Score
}
