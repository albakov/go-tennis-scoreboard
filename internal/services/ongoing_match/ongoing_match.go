package ongoing_match

import (
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/match"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/score"
)

type OngoingMatchService struct {
	storageMatch match.StorageMatch
	storageScore score.StorageScore
}

func New(config *config.Config) *OngoingMatchService {
	return &OngoingMatchService{
		storageMatch: match.New(config.PathToDB),
		storageScore: score.New(config.PathToDB),
	}
}

func (m *OngoingMatchService) OngoingMatch(uuid string, currentPlayerID int64) (entity.OngoingMatch, error) {
	c := entity.OngoingMatch{}

	var err error

	c.MatchGame, err = m.storageMatch.ByUUID(uuid)
	if err != nil {
		return c, err
	}

	if c.MatchGame.PlayerOneID != currentPlayerID && c.MatchGame.PlayerTwoID != currentPlayerID {
		return c, err
	}

	c.ScoreCurrentPlayer, err = m.scorePlayer(c.MatchGame.ID, currentPlayerID)
	if err != nil {
		return c, err
	}

	c.ScoreAnotherPlayer, err = m.scorePlayer(c.MatchGame.ID, m.anotherPlayerID(c.MatchGame, currentPlayerID))
	if err != nil {
		return c, err
	}

	return c, nil
}

func (m *OngoingMatchService) scorePlayer(matchID, playerID int64) (entity.Score, error) {
	scoreCurrentPlayer, err := m.storageScore.ByMatchIDAndPlayerID(matchID, playerID)
	if err != nil {
		return entity.Score{}, err
	}

	return scoreCurrentPlayer, nil
}

func (m *OngoingMatchService) anotherPlayerID(matchGame entity.Match, currentPlayerID int64) int64 {
	if matchGame.PlayerOneID == currentPlayerID {
		return matchGame.PlayerTwoID
	}

	return matchGame.PlayerOneID
}
