package models

type Game struct {
	Id            int
	Date          string
	Time          string
	Duration      string
	Place         string
	Level         string
	Players       []string
	NumberOfSpots int
	CreatorId     int64
	IsPublished   bool
}

func (g Game) IsFull() bool {
	return len(g.Players) == g.NumberOfSpots
}
