package score_calculate

import (
	"database/sql"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"testing"
)

func TestCalculate(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     0,
			Points:    0,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     0,
			Points:    0,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 15 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 15)
	}
}

func TestAdvantage(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if !obj.OngoingMatch.ScoreCurrentPlayer.Advantage {
		t.Errorf("Result was incorrect, got: %t, want: %t.", obj.OngoingMatch.ScoreCurrentPlayer.Advantage, true)
	}
}

func TestResetAdvantage(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: true,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Advantage {
		t.Errorf("Result was incorrect, got: %t, want: %t.", !obj.OngoingMatch.ScoreCurrentPlayer.Advantage, false)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Advantage {
		t.Errorf("Result was incorrect, got: %t, want: %t.", !obj.OngoingMatch.ScoreAnotherPlayer.Advantage, false)
	}
}

func TestWinWithAdvantage(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: true,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     0,
			Points:    40,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 0)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Games != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Games, 1)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Points != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Points, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Games != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Games, 0)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Advantage {
		t.Errorf("Result was incorrect, got: %t, want: %t.", !obj.OngoingMatch.ScoreCurrentPlayer.Advantage, false)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Advantage {
		t.Errorf("Result was incorrect, got: %t, want: %t.", !obj.OngoingMatch.ScoreAnotherPlayer.Advantage, false)
	}
}

func TestTieBreak(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     6,
			Points:    0,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     6,
			Points:    0,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if !obj.isTieBreak() {
		t.Errorf("Result was incorrect, got: %t, want: %t.", obj.isTieBreak(), true)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 1)
	}
}

func TestNotWinTieBreak(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     6,
			Points:    7,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     6,
			Points:    6,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 8 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 8)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Sets != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Sets, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Sets != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Sets, 0)
	}
}

func TestWinTieBreak(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      0,
			Games:     6,
			Points:    7,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      0,
			Games:     6,
			Points:    5,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 0)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Sets != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Sets, 1)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Games != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Games, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Games != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Games, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Sets != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Sets, 0)
	}
}

func TestNotWinYetMatch(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      1,
			Games:     5,
			Points:    40,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      1,
			Games:     5,
			Points:    30,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 0)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Sets != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Sets, 1)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Games != 6 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Games, 6)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Games != 5 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Games, 5)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Sets != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Sets, 1)
	}
}

func TestWinMatch(t *testing.T) {
	item := entity.OngoingMatch{
		MatchGame: entity.Match{
			ID:          1,
			UUID:        "1",
			PlayerOneID: 1,
			PlayerTwoID: 2,
			WinnerID:    sql.NullInt64{},
		},
		ScoreCurrentPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  1,
			Sets:      1,
			Games:     5,
			Points:    40,
			Advantage: false,
		},
		ScoreAnotherPlayer: entity.Score{
			MatchID:   1,
			PlayerID:  2,
			Sets:      1,
			Games:     4,
			Points:    30,
			Advantage: false,
		},
	}
	obj := New(&item)
	obj.Calculate()

	if obj.OngoingMatch.ScoreCurrentPlayer.Points != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Points, 0)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Sets != 2 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Sets, 2)
	}

	if obj.OngoingMatch.ScoreCurrentPlayer.Games != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreCurrentPlayer.Games, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Games != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Games, 0)
	}

	if obj.OngoingMatch.ScoreAnotherPlayer.Sets != 1 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", obj.OngoingMatch.ScoreAnotherPlayer.Sets, 1)
	}
}
