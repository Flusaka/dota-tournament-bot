package types

type Match struct {
	ID            int16  `json:"id"`
	Radiant       *Team  `json:"radiant"`
	Dire          *Team  `json:"dire"`
	ScheduledTime int64  `json:"scheduledTime"`
	StreamUrl     string `json:"streamUrl"`
}

func NewMatch(id int16, radiant *Team, dire *Team, scheduledTime int64, streamUrl string) *Match {
	return &Match{
		id, radiant, dire, scheduledTime, streamUrl,
	}
}
