package score

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
)

const f = "storage.Score"

type StorageScore interface {
	Create(matchID, playerID int64) (entity.Score, error)
	ByMatchIDAndPlayerID(matchID, playerID int64) (entity.Score, error)
	Update(item entity.Score) error
	DeleteByMatchID(matchID int64) error
}

type Score struct {
	pathToDb string
}

func New(pathToDb string) *Score {
	return &Score{
		pathToDb: pathToDb,
	}
}

func (c *Score) Create(matchID, playerID int64) (entity.Score, error) {
	const op = "Create"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	stmt, err := db.Prepare(
		"INSERT INTO Scores (MatchID, PlayerID, Sets, Games, Points, Advantage) VALUES (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		util.LogError(f, op, err)

		return entity.Score{}, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	_, err = stmt.Exec(matchID, playerID, 0, 0, 0, 0)
	if err != nil {
		util.LogError(f, op, err)

		return entity.Score{}, err
	}

	return c.ByMatchIDAndPlayerID(matchID, playerID)
}

func (c *Score) ByMatchIDAndPlayerID(matchID, playerID int64) (entity.Score, error) {
	const op = "ByMatchIDAndPlayerID"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	item := entity.Score{}

	row := db.QueryRow(
		"SELECT MatchID, PlayerID, Sets, Games, Points, Advantage FROM Scores WHERE MatchID = ? AND PlayerID = ?",
		matchID,
		playerID,
	)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return item, row.Err()
	}

	err = row.Scan(&item.MatchID, &item.PlayerID, &item.Sets, &item.Games, &item.Points, &item.Advantage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Score{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Score{}, err
	}

	return item, nil
}

func (c *Score) Update(item entity.Score) error {
	const op = "Update"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	stmt, err := db.Prepare(
		"UPDATE Scores SET Sets = ?, Games = ?, Points = ?, Advantage = ? WHERE MatchID = ? AND PlayerID = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Sets, item.Games, item.Points, item.Advantage, item.MatchID, item.PlayerID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Score) DeleteByMatchID(matchID int64) error {
	const op = "DeleteByMatchID"

	db, err := sql.Open("sqlite3", c.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(db)

	stmt, err := db.Prepare("DELETE FROM Scores WHERE MatchID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(matchID)
	if err != nil {
		return err
	}

	return nil
}
