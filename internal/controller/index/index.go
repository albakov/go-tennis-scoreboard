package index

import (
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/player"
	"net/http"
)

type Controller struct {
	storageCurrencies player.StoragePlayer
	commonController  controller.ServerResponse
}

func New(commonController controller.ServerResponse) *Controller {
	return &Controller{commonController: commonController}
}

func (cc *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		cc.commonController.ShowError(
			w,
			http.StatusNotFound,
			controller.ErrorMessage{Title: "Not found", Code: http.StatusNotFound},
		)

		return
	}

	if r.Method != http.MethodGet {
		cc.commonController.ShowMethodNotAllowedError(w)

		return
	}

	cc.commonController.ShowResponse(w, http.StatusOK, "index.html", controller.PageData{PageTitle: "Home"})
}
