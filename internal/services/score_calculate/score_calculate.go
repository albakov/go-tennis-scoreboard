package score_calculate

import (
	"github.com/albakov/go-tennis-scoreboard/internal/entity"
)

type ScoreCalculateService struct {
	OngoingMatch *entity.OngoingMatch
}

func New(Calculated *entity.OngoingMatch) *ScoreCalculateService {
	return &ScoreCalculateService{OngoingMatch: Calculated}
}

func (s *ScoreCalculateService) Calculate() {
	if s.isTieBreak() {
		s.calculateTieBreak()

		return
	}

	if s.OngoingMatch.ScoreCurrentPlayer.Advantage {
		s.calculateIfAdvantage()

		return
	}

	s.calculatePoints()
}

func (s *ScoreCalculateService) isTieBreak() bool {
	return s.OngoingMatch.ScoreCurrentPlayer.Games == 6 && s.OngoingMatch.ScoreAnotherPlayer.Games == 6
}

func (s *ScoreCalculateService) calculateIfAdvantage() {
	s.resetAdvantage()
	s.resetPoints()
	s.calculateGames()
}

func (s *ScoreCalculateService) calculatePoints() {
	if s.OngoingMatch.ScoreCurrentPlayer.Points == 40 {
		if s.OngoingMatch.ScoreAnotherPlayer.Points == 40 {
			if s.OngoingMatch.ScoreAnotherPlayer.Advantage {
				s.resetAdvantage()
			} else {
				s.OngoingMatch.ScoreCurrentPlayer.Advantage = true
			}

			return
		}

		s.resetPoints()
		s.calculateGames()

		return
	}

	s.addPointsDefault()
}

func (s *ScoreCalculateService) calculateGames() {
	s.OngoingMatch.ScoreCurrentPlayer.Games += 1

	if s.OngoingMatch.ScoreCurrentPlayer.Games >= 6 {
		if s.isGoingToWinSet() {
			s.resetPoints()
			s.resetGames()
			s.OngoingMatch.ScoreCurrentPlayer.Sets += 1
		}
	}
}

func (s *ScoreCalculateService) calculateTieBreak() {
	if s.OngoingMatch.ScoreCurrentPlayer.Points >= 7 {
		if s.isGoingToWinSetAfterTieBreak() {
			s.resetPoints()
			s.resetGames()
			s.OngoingMatch.ScoreCurrentPlayer.Sets += 1

			return
		}
	}

	s.OngoingMatch.ScoreCurrentPlayer.Points += 1
}

func (s *ScoreCalculateService) resetPoints() {
	s.OngoingMatch.ScoreCurrentPlayer.Points = 0
	s.OngoingMatch.ScoreAnotherPlayer.Points = 0
}

func (s *ScoreCalculateService) resetAdvantage() {
	s.OngoingMatch.ScoreCurrentPlayer.Advantage = false
	s.OngoingMatch.ScoreAnotherPlayer.Advantage = false
}

func (s *ScoreCalculateService) resetGames() {
	s.OngoingMatch.ScoreCurrentPlayer.Games = 0
	s.OngoingMatch.ScoreAnotherPlayer.Games = 0
}

func (s *ScoreCalculateService) addPointsDefault() {
	if s.OngoingMatch.ScoreCurrentPlayer.Points == 30 {
		s.OngoingMatch.ScoreCurrentPlayer.Points = 40
	} else {
		s.OngoingMatch.ScoreCurrentPlayer.Points += 15
	}
}

func (s *ScoreCalculateService) isGoingToWinSet() bool {
	return s.OngoingMatch.ScoreCurrentPlayer.Games-s.OngoingMatch.ScoreAnotherPlayer.Games >= 2
}

func (s *ScoreCalculateService) isGoingToWinSetAfterTieBreak() bool {
	return s.OngoingMatch.ScoreCurrentPlayer.Points-s.OngoingMatch.ScoreAnotherPlayer.Points >= 2
}
