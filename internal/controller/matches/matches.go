package matches

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/controller"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/match"
	"net/http"
	"strconv"
	"strings"
)

type Controller struct {
	config           *config.Config
	commonController controller.ServerResponse
	storageMatch     match.StorageMatch
}

type PageMatches struct {
	Matches []entity.MatchWithPlayer
	Paginator
	Filter entity.Filter
}

type Paginator struct {
	CanShowPrevPage, CanShowNextPage bool
	CurrentPage, LastPage            int64
	PrevPageUrl, NextPageUrl         string
	Pages                            []Page
}

type Page struct {
	Number int64
	Url    string
}

func New(config *config.Config, commonController controller.ServerResponse) *Controller {
	return &Controller{
		config:           config,
		commonController: commonController,
		storageMatch:     match.New(config.PathToDB),
	}
}

func (cc *Controller) MatchesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cc.commonController.ShowMethodNotAllowedError(w)

		return
	}

	filter := entity.Filter{Column: "filter_by_player_name"}
	filter.Value = r.URL.Query().Get("filter_by_player_name")

	pages, err := cc.storageMatch.Pages(cc.config.Game.MatchesPerPage, filter)
	if err != nil {
		cc.commonController.ShowServerError(w)

		return
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 || page > pages {
		page = 1
	}

	cc.commonController.ShowResponse(
		w,
		http.StatusOK,
		"matches.html",
		controller.PageData{
			PageTitle: "Finished Matches",
			Data: PageMatches{
				Matches: cc.storageMatch.Paginate(
					cc.config.Game.MatchesPerPage,
					(page-1)*cc.config.Game.MatchesPerPage,
					filter,
				),
				Filter: filter,
				Paginator: Paginator{
					CurrentPage:     page,
					LastPage:        pages,
					CanShowPrevPage: page > 1,
					CanShowNextPage: page < pages,
					PrevPageUrl:     cc.url(page-1, filter),
					NextPageUrl:     cc.url(page+1, filter),
					Pages:           cc.matchPages(pages, filter),
				},
			},
		},
	)
}

func (cc *Controller) matchPages(pages int64, filter entity.Filter) []Page {
	s := make([]Page, pages)

	var i int64 = 0

	for ; i < pages; i++ {
		page := i + 1
		s[i].Number = page
		s[i].Url = cc.url(page, filter)
	}

	return s
}

func (cc *Controller) url(page int64, filter entity.Filter) string {
	params := []string{}

	if page > 1 {
		params = append(params, fmt.Sprintf("page=%d", page))
	}

	if filter.Value != "" {
		params = append(params, fmt.Sprintf("%s=%s", filter.Column, filter.Value))
	}

	query := ""

	if len(params) > 0 {
		query = fmt.Sprintf("?%s", strings.Join(params, "&"))
	}

	return fmt.Sprintf("/matches%s", query)
}
