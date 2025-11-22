package dto

// TideMonth represents the tide information for a given calendar month.
type TideMonth struct {
	Month string    `json:"month"`
	Days  []TideDay `json:"days"`
}

// TideDay holds the list of tide observations for a specific day.
type TideDay struct {
	Day          int               `json:"day"`
	Observations []TideObservation `json:"observations"`
}

// TideObservation captures the time, height and high/low indicator of a tide event.
type TideObservation struct {
	Time   string  `json:"time"`
	Height float64 `json:"height"`
	Level  string  `json:"level"`
}
