package jellyfindata

import "time"

type UserData struct {
	LastPlayedDate time.Time `json:"LastPlayedDate"`
}

type Item struct {
	UserData     UserData `json:"UserData"`
	ID           string   `json:"Id"`
	SeriesName   string   `json:"SeriesName,omitempty"`
	SeriesID     string   `json:"SeriesId,omitempty"`
	RunTimeTicks int64    `json:"RunTimeTicks"`
	Type         string   `json:"Type"`
	Name         string   `json:"Name"`
}

type ActivityResponse struct {
	Items []Item `json:"Items"`
}

type AggregatedSeries struct {
	ID    string
	Name  string
	Ticks int64
}

type AggregatedMovie struct {
	ID    string
	Name  string
	Ticks int64
}
