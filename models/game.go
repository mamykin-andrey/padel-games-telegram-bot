package models

import "fmt"

type Game struct {
	Id            int
	Date          string
	Time          string
	Duration      string
	Place         string
	Level         string
	Players       []string
	NumberOfSpots int
	IsPublished   bool
}

func (g Game) String() string {
	return fmt.Sprint("Id: ", g.Id, ", date: ", g.Date, ", time: ", g.Time, ", duration: ", g.Duration, ", place: ", g.Place, ", level: ", g.Level, ", players: ", g.Players)
}

func (g Game) IsFull() bool {
	return len(g.Players) == g.NumberOfSpots
}
