package types

type Match struct {
	Radiant       *Team
	Dire          *Team
	ScheduledTime int64
	StreamUrl     string
}

func NewMatch(radiant *Team, dire *Team, scheduledTime int64, streamUrl string) *Match {
	return &Match{
		radiant, dire, scheduledTime, streamUrl,
	}
}
