package main

import (
	"github.com/albakov/go-tennis-scoreboard/dbinit"
	"github.com/albakov/go-tennis-scoreboard/internal/app"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
)

func main() {
	c := config.MustNew()
	dbinit.New(c).MustCreateDatabaseIfNotExists()
	app.New(c).MustStart()
}
