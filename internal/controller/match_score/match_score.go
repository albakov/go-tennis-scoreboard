package match_score

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/services/match_persistence"
	"github.com/albakov/go-tennis-scoreboard/internal/services/ongoing_match"
	"github.com/albakov/go-tennis-scoreboard/internal/services/score_calculate"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/match"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/player"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/score"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
	"net/http"
	"strconv"
)

const f = "controller.match_score"

type Controller struct {
	config           *config.Config
	commonController controller.ServerResponse
	storagePlayer    player.StoragePlayer
	storageMatch     match.StorageMatch
	storageScore     score.StorageScore
}

func New(config *config.Config, commonController controller.ServerResponse) *Controller {
	return &Controller{
		config:           config,
		commonController: commonController,
		storagePlayer:    player.New(config.PathToDB),
		storageMatch:     match.New(config.PathToDB),
		storageScore:     score.New(config.PathToDB),
	}
}

func (cc *Controller) MatchScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.renderScorePage(w, r)

		return
	}

	if r.Method == http.MethodPost {
		cc.handleScoreRequest(w, r)

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w)
}

func (cc *Controller) renderScorePage(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		cc.commonController.ShowNotFound(w)

		return
	}

	matchGame, err := cc.storageMatch.ByUUID(uuid)
	if err != nil {
		cc.commonController.ShowNotFound(w)

		return
	}

	if matchGame.IsFinished() {
		cc.renderFinishedPage(w, matchGame)

		return
	}

	matchScore, err := cc.matchScore(matchGame)
	if err != nil {
		cc.commonController.ShowNotFound(w)

		return
	}

	cc.commonController.ShowResponse(w, http.StatusOK, "match-score.html", controller.PageData{
		PageTitle: "Match Score",
		Data:      matchScore,
	})
}

func (cc *Controller) handleScoreRequest(w http.ResponseWriter, r *http.Request) {
	v := r.FormValue("player")
	currentPlayerID, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		cc.commonController.ShowServerError(w)

		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		cc.commonController.ShowNotFound(w)

		return
	}

	ongoingMatchService := ongoing_match.New(cc.config)
	ongoingMatch, err := ongoingMatchService.OngoingMatch(uuid, currentPlayerID)
	if err != nil {
		cc.commonController.ShowNotFound(w)

		return
	}

	if ongoingMatch.MatchGame.IsFinished() {
		http.Redirect(w, r, fmt.Sprintf("/match-score?uuid=%s", uuid), http.StatusFound)

		return
	}

	score_calculate.New(&ongoingMatch).Calculate()
	err = match_persistence.New(cc.config).Sync(ongoingMatch)
	if err != nil {
		cc.commonController.ShowServerError(w)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/match-score?uuid=%s", uuid), http.StatusFound)
}

func (cc *Controller) matchScore(matchGame entity.Match) (entity.MatchScore, error) {
	const op = "matchScore"

	matchScore := entity.MatchScore{MatchGame: matchGame}

	var err error

	matchScore.ScorePlayerOne, err = cc.storageScore.ByMatchIDAndPlayerID(matchGame.ID, matchGame.PlayerOneID)
	if err != nil {
		util.LogError(f, op, err)

		return matchScore, err
	}

	matchScore.ScorePlayerTwo, err = cc.storageScore.ByMatchIDAndPlayerID(matchGame.ID, matchGame.PlayerTwoID)
	if err != nil {
		util.LogError(f, op, err)

		return matchScore, err
	}

	matchScore.PlayerOne, err = cc.storagePlayer.ByID(matchScore.ScorePlayerOne.PlayerID)
	if err != nil {
		util.LogError(f, op, err)

		return matchScore, err
	}

	matchScore.PlayerTwo, err = cc.storagePlayer.ByID(matchScore.ScorePlayerTwo.PlayerID)
	if err != nil {
		util.LogError(f, op, err)

		return matchScore, err
	}

	return matchScore, nil
}

func (cc *Controller) renderFinishedPage(w http.ResponseWriter, matchGame entity.Match) {
	winner, err := cc.storagePlayer.ByID(matchGame.WinnerID.Int64)
	if err != nil {
		return
	}

	loserID := matchGame.PlayerTwoID

	if winner.ID == matchGame.PlayerTwoID {
		loserID = matchGame.PlayerOneID
	}

	loser, err := cc.storagePlayer.ByID(loserID)

	type MatchFinished struct {
		Winner entity.Player
		Loser  entity.Player
	}

	cc.commonController.ShowResponse(w, http.StatusOK, "match-finished.html", controller.PageData{
		PageTitle: "Match Score",
		Data: MatchFinished{
			Winner: winner,
			Loser:  loser,
		},
	})
}
