package match_persistence

import (
	"github.com/albakov/go-tennis-scoreboard/internal/config"
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/match"
	"github.com/albakov/go-tennis-scoreboard/internal/storage/score"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
)

const f = "match_persistence.MatchPersistence"

type MatchPersistence struct {
	config       *config.Config
	storageMatch match.StorageMatch
	storageScore score.StorageScore
}

func New(config *config.Config) *MatchPersistence {
	return &MatchPersistence{
		config:       config,
		storageMatch: match.New(config.PathToDB),
		storageScore: score.New(config.PathToDB),
	}
}

func (mp *MatchPersistence) Sync(calculated entity.OngoingMatch) error {
	const op = "Sync"

	err := mp.storageScore.Update(calculated.ScoreCurrentPlayer)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	err = mp.storageScore.Update(calculated.ScoreAnotherPlayer)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	if calculated.ScoreCurrentPlayer.Sets == mp.config.Game.MaxSetsToWin {
		err := mp.storageScore.DeleteByMatchID(calculated.MatchGame.ID)
		if err != nil {
			util.LogError(f, op, err)

			return err
		}

		err = mp.storageMatch.SetWinnerID(calculated.MatchGame.ID, calculated.ScoreCurrentPlayer.PlayerID)
		if err != nil {
			util.LogError(f, op, err)

			return err
		}
	}

	return nil
}
