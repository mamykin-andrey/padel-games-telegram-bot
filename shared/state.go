package shared

import (
	"main/models"
	"sync"
)

type GamesState struct {
	mutex sync.Mutex
	games []models.Game
}

func (s *GamesState) Games() []models.Game {
	return s.games
}

func (s *GamesState) Add(game models.Game) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.games = append(s.games, game)
}

func (s *GamesState) Remove(index int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.games = append(s.games[:index], s.games[index+1:]...)
}

// TODO: Add a separate state for games drafts
var State GamesState
