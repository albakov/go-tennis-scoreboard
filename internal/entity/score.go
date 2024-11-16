package entity

type Score struct {
	MatchID   int64
	PlayerID  int64
	Sets      int64
	Games     int64
	Points    int64
	Advantage bool
}
