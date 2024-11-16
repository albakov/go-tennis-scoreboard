package app

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"github.com/albakov/go-tennis-scoreboard/internal/controller/index"
	"github.com/albakov/go-tennis-scoreboard/internal/controller/match_score"
	"github.com/albakov/go-tennis-scoreboard/internal/controller/matches"
	"github.com/albakov/go-tennis-scoreboard/internal/controller/new_match"
	"net/http"
	"strings"
)

type App struct {
	mux                  *http.ServeMux
	config               *config.Config
	indexController      *index.Controller
	newMatchController   *new_match.Controller
	matchesController    *matches.Controller
	matchScoreController *match_score.Controller
}

func New(config *config.Config) *App {
	commonController := controller.New()

	return &App{
		mux:                  http.NewServeMux(),
		config:               config,
		indexController:      index.New(commonController),
		newMatchController:   new_match.New(config, commonController),
		matchesController:    matches.New(config, commonController),
		matchScoreController: match_score.New(config, commonController),
	}
}

func (a *App) MustStart() {
	a.setRoutes()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.config.Host, a.config.Port), a)
	if err != nil {
		panic(err)
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) setRoutes() {
	a.setAssets()
	a.mux.HandleFunc("/", a.indexController.IndexHandler)
	a.mux.HandleFunc("/new-match", a.newMatchController.NewMatchHandler)
	a.mux.HandleFunc("/matches", a.matchesController.MatchesHandler)
	a.mux.HandleFunc("/match-score", a.matchScoreController.MatchScoreHandler)
}

func (a *App) setAssets() {
	a.mux.Handle("/css/", a.noDirListing(http.FileServer(http.Dir("./view"))))
	a.mux.Handle("/js/", a.noDirListing(http.FileServer(http.Dir("./view"))))
	a.mux.Handle("/images/", a.noDirListing(http.FileServer(http.Dir("./view"))))
}

func (a *App) noDirListing(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/404", http.StatusFound)

			return
		}

		h.ServeHTTP(w, r)
	}
}
