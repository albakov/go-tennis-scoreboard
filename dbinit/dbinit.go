package dbinit

import (
	"database/sql"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type DBInit struct {
	pathToDb string
}

func New(config *config.Config) *DBInit {
	return &DBInit{
		pathToDb: config.PathToDB,
	}
}

func (d *DBInit) MustCreateDatabaseIfNotExists() {
	if !d.mustCheckIsDatabaseExists() {
		_, err := os.Create(d.pathToDb)
		if err != nil {
			panic(err)
		}
	}

	d.mustCreateTablesIfNotExists()
}

func (d *DBInit) mustCheckIsDatabaseExists() bool {
	_, err := os.Stat(d.pathToDb)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}

		panic(err)
	}

	return true
}

func (d *DBInit) mustCreateTablesIfNotExists() {
	db, err := sql.Open("sqlite3", d.pathToDb)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Players (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    	FullName VARCHAR(255) NOT NULL)`,
	)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Matches (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    	UUID VARCHAR(36) NOT NULL, 
    	PlayerOneID INT NOT NULL, 
    	PlayerTwoID INT NOT NULL, 
    	WinnerID INT DEFAULT NULL NULL, 
    	FOREIGN KEY (PlayerOneID) REFERENCES Players (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	FOREIGN KEY (PlayerTwoID) REFERENCES Players (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	FOREIGN KEY (WinnerID) REFERENCES Players (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	UNIQUE(UUID) ON CONFLICT ABORT)`,
	)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Scores (
    	MatchID INT NOT NULL, 
    	PlayerID INT NOT NULL, 
    	Sets INT DEFAULT 0 NOT NULL, 
    	Games INT DEFAULT 0 NOT NULL, 
    	Points INT DEFAULT 0 NOT NULL, 
    	Advantage INT DEFAULT 0 NOT NULL, 
    	FOREIGN KEY (MatchID) REFERENCES Matches (ID) ON DELETE CASCADE ON UPDATE NO ACTION,
    	FOREIGN KEY (PlayerID) REFERENCES Players (ID) ON DELETE CASCADE ON UPDATE NO ACTION)`,
	)
	if err != nil {
		panic(err)
	}
}
