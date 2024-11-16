package entity

type MatchScore struct {
	MatchGame      Match
	ScorePlayerOne Score
	ScorePlayerTwo Score
	PlayerOne      Player
	PlayerTwo      Player
}
