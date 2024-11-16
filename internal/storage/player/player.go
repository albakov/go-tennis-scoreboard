package player

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
	"strings"
)

const f = "storage.Player"

type StoragePlayer interface {
	ByID(ID int64) (entity.Player, error)
	ByName(name string) (entity.Player, error)
	FirstOrCreate(name string) (entity.Player, error)
}

type Player struct {
	pathToDb string
}

func New(pathToDb string) *Player {
	return &Player{
		pathToDb: pathToDb,
	}
}

func (c *Player) ByID(ID int64) (entity.Player, error) {
	const op = "ByID"

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

	item := entity.Player{}

	row := db.QueryRow("SELECT ID, FullName FROM Players WHERE ID = ?", ID)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return item, row.Err()
	}

	err = row.Scan(&item.ID, &item.FullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Player{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Player{}, err
	}

	return item, nil
}

func (c *Player) ByName(name string) (entity.Player, error) {
	const op = "ByName"

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

	item := entity.Player{}

	row := db.QueryRow("SELECT ID, FullName FROM Players WHERE FullName = ?", name)
	if row.Err() != nil {
		util.LogError(f, op, row.Err())

		return item, row.Err()
	}

	err = row.Scan(&item.ID, &item.FullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Player{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Player{}, err
	}

	return item, nil
}

func (c *Player) FirstOrCreate(name string) (entity.Player, error) {
	const op = "FirstOrCreate"

	item, err := c.ByName(name)
	if err != nil {
		if !errors.Is(err, storage.EntitiesNotFoundError) {
			return entity.Player{}, err
		}
	}

	if item.ID != 0 {
		return item, nil
	}

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

	stmt, err := db.Prepare("INSERT INTO Players (FullName) VALUES (?)")
	if err != nil {
		util.LogError(f, op, err)

		return entity.Player{}, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			util.LogError(f, op, err)
		}
	}(stmt)

	_, err = stmt.Exec(name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return entity.Player{}, storage.EntityAlreadyExistsError
		}

		util.LogError(f, op, err)

		return entity.Player{}, err
	}

	return c.ByName(name)
}
