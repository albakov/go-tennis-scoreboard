package match

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
	uuid2 "github.com/google/uuid"
	"math"
	"strings"
)

const f = "storage.Match"

type StorageMatch interface {
	Create(playerOneID, playerTwoID int64) (entity.Match, error)
	ByUUID(uuid string) (entity.Match, error)
	SetWinnerID(id int64, playerID int64) error
	Pages(limit int64, filter entity.Filter) (int64, error)
	Paginate(limit int64, perPage int64, filter entity.Filter) []entity.MatchWithPlayer
}

type Match struct {
	pathToDb string
}

func New(pathToDb string) *Match {
	return &Match{
		pathToDb: pathToDb,
	}
}

func (c *Match) Create(playerOneID, playerTwoID int64) (entity.Match, error) {
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

	stmt, err := db.Prepare("INSERT INTO Matches (UUID, PlayerOneID, PlayerTwoID, WinnerID) VALUES (?, ?, ?, null)")
	if err != nil {
		util.LogError(f, op, err)

		return entity.Match{}, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	uuid := uuid2.Must(uuid2.NewRandom()).String()

	_, err = stmt.Exec(uuid, playerOneID, playerTwoID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return entity.Match{}, storage.EntityAlreadyExistsError
		}

		util.LogError(f, op, err)

		return entity.Match{}, err
	}

	return c.ByUUID(uuid)
}

func (c *Match) ByUUID(uuid string) (entity.Match, error) {
	const op = "ByUUID"

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

	item := entity.Match{}

	row := db.QueryRow("SELECT ID, UUID, PlayerOneID, PlayerTwoID, WinnerID FROM Matches WHERE UUID = ?", uuid)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return item, row.Err()
	}

	err = row.Scan(&item.ID, &item.UUID, &item.PlayerOneID, &item.PlayerTwoID, &item.WinnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Match{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Match{}, err
	}

	return item, nil
}

func (c *Match) SetWinnerID(id int64, winnerID int64) error {
	const op = "SetWinnerID"

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

	stmt, err := db.Prepare("UPDATE Matches SET WinnerID = ? WHERE ID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(winnerID, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Match) Pages(limit int64, filter entity.Filter) (int64, error) {
	const op = "Pages"

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

	var count int64 = 0
	args, conditions := c.argsForPaginator(filter)

	row := db.QueryRow(
		fmt.Sprintf(
			`SELECT count(Matches.ID) AS count FROM Matches 
                LEFT JOIN Players AS P1 ON P1.ID = Matches.PlayerOneID
                LEFT JOIN Players AS P2 ON P2.ID = Matches.PlayerTwoID
				WHERE Matches.WinnerID IS NOT NULL %s ORDER BY Matches.ID ASC`,
			conditions,
		),
		args...,
	)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return count, row.Err()
	}

	err = row.Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	return int64(math.Ceil(float64(count) / float64(limit))), nil
}

func (c *Match) Paginate(limit int64, offset int64, filter entity.Filter) []entity.MatchWithPlayer {
	const op = "Paginate"

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

	items := []entity.MatchWithPlayer{}
	args, conditions := c.argsForPaginator(filter)
	args = append(args, offset, limit)

	stmt, err := db.Query(
		fmt.Sprintf(
			`SELECT P1.FullName, P2.FullName, W.FullName FROM Matches 
                LEFT JOIN Players AS P1 ON P1.ID = Matches.PlayerOneID
                LEFT JOIN Players AS P2 ON P2.ID = Matches.PlayerTwoID
                LEFT JOIN Players AS W ON W.ID = Matches.WinnerID
        		WHERE Matches.WinnerID IS NOT NULL %s ORDER BY Matches.ID ASC LIMIT ?,?`,
			conditions,
		),
		args...,
	)
	if err != nil {
		util.LogError(f, op, err)

		return items
	}
	defer stmt.Close()

	for stmt.Next() {
		item := entity.MatchWithPlayer{}

		err = stmt.Scan(&item.PlayerOneFullName, &item.PlayerTwoFullName, &item.WinnerFullName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []entity.MatchWithPlayer{}
			}

			util.LogError(f, op, err)

			return []entity.MatchWithPlayer{}
		}

		items = append(items, item)
	}

	return items
}

func (c *Match) argsForPaginator(filter entity.Filter) ([]any, string) {
	var args []any
	conditions := ""

	if filter.Value != "" {
		conditions += " AND (P1.FullName LIKE ? OR P2.FullName LIKE ?)"
		value := fmt.Sprintf("%%%s%%", filter.Value)
		args = append(args, value, value)
	}

	return args, conditions
}
