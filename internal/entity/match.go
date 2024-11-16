package entity

import "database/sql"

type Match struct {
	ID          int64
	UUID        string
	PlayerOneID int64
	PlayerTwoID int64
	WinnerID    sql.NullInt64
}

func (m *Match) IsFinished() bool {
	return m.WinnerID.Valid
}

type MatchWithPlayer struct {
	PlayerOneFullName, PlayerTwoFullName, WinnerFullName string
}
