package types

type Match struct {
	Radiant       *Team  `json:"radiant"`
	Dire          *Team  `json:"dire"`
	ScheduledTime int64  `json:"scheduledTime"`
	StreamUrl     string `json:"streamUrl"`
}

func NewMatch(radiant *Team, dire *Team, scheduledTime int64, streamUrl string) *Match {
	return &Match{
		radiant, dire, scheduledTime, streamUrl,
	}
}
