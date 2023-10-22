package bot

import (
	"github.com/flusaka/dota-tournament-bot/datasource/types"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

type MatchStartedNotification struct {
	Users []string
	Match *types.Match
}

type MatchEventNotifier struct {
	startedNotifications map[int16]*MatchStartedNotification
	MatchStarted         chan *MatchStartedNotification
	mux                  sync.Mutex
	cancel               <-chan bool
}

func NewMatchEventNotifier(cancel <-chan bool) *MatchEventNotifier {
	matchEventNotifier := &MatchEventNotifier{
		startedNotifications: make(map[int16]*MatchStartedNotification),
		MatchStarted:         make(chan *MatchStartedNotification, 5),
		cancel:               cancel,
	}

	return matchEventNotifier
}

func (r *MatchEventNotifier) AddUserToNotificationsForMatch(match *types.Match, userID string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if existing, ok := r.startedNotifications[match.ID]; ok {
		if !slices.Contains(existing.Users, userID) {
			existing.Users = append(existing.Users, userID)
		}
	} else {
		notification := &MatchStartedNotification{
			Users: []string{userID},
			Match: match,
		}

		r.startedNotifications[match.ID] = notification
		r.startMatchTicker(notification)
	}
}

func (r *MatchEventNotifier) startMatchTicker(notification *MatchStartedNotification) {
	go func() {
		now := time.Now().UTC()
		matchStart := time.Unix(notification.Match.ScheduledTime, 0).UTC()
		duration := matchStart.Sub(now)
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				{

					r.mux.Lock()
					delete(r.startedNotifications, notification.Match.ID)
					r.mux.Unlock()

					r.MatchStarted <- notification
					return
				}
			case <-r.cancel:
				{
					return
				}
			}
		}
	}()
}
