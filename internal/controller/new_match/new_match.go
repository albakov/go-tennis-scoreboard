package new_match

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/match"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/player"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/score"
	matchValidator "github.com/albakov/go-tennis-scoreboard/internal/validation/new_match"
	"html/template"
	"net/http"
)

type Controller struct {
	commonController controller.ServerResponse
	storagePlayer    player.StoragePlayer
	storageMatch     match.StorageMatch
	storageScore     score.StorageScore
}

func New(config *config.Config, commonController controller.ServerResponse) *Controller {
	return &Controller{
		commonController: commonController,
		storagePlayer:    player.New(config.PathToDB),
		storageMatch:     match.New(config.PathToDB),
		storageScore:     score.New(config.PathToDB),
	}
}

func (cc *Controller) NewMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.commonController.ShowResponse(
			w,
			http.StatusOK,
			"new-match.html",
			controller.PageData{PageTitle: "New Match"},
		)

		return
	}

	if r.Method == http.MethodPost {
		cc.createNewMatch(w, r)

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w)
}

func (cc *Controller) createNewMatch(w http.ResponseWriter, r *http.Request) {
	validator := matchValidator.NewMatch(r, map[string]string{"playerOne": "", "playerTwo": ""})
	validator.Validate()

	if !validator.IsValid() {
		cc.showError(w, validator.ErrorMessage())

		return
	}

	playerOne, err := cc.storagePlayer.FirstOrCreate(validator.Field("playerOne"))
	if err != nil {
		cc.showError(w, controller.MessageServerError)

		return
	}

	playerTwo, err := cc.storagePlayer.FirstOrCreate(validator.Field("playerTwo"))
	if err != nil {
		cc.showError(w, controller.MessageServerError)

		return
	}

	matchGame, err := cc.storageMatch.Create(playerOne.ID, playerTwo.ID)
	if err != nil {
		cc.showError(w, controller.MessageServerError)

		return
	}

	_, err = cc.storageScore.Create(matchGame.ID, playerOne.ID)
	if err != nil {
		cc.showError(w, controller.MessageServerError)

		return
	}

	_, err = cc.storageScore.Create(matchGame.ID, playerTwo.ID)
	if err != nil {
		cc.showError(w, controller.MessageServerError)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/match-score?uuid=%s", matchGame.UUID), http.StatusFound)
}

func (cc *Controller) showError(w http.ResponseWriter, message string) {
	cc.commonController.BackWithError(
		w,
		"new-match.html",
		controller.PageData{PageTitle: "New Match", ErrorMessage: template.HTML(message)},
	)
}
