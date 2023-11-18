package bot

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"golang.org/x/exp/slices"
	"log"
	"sync"
	"time"
)

type MatchStartedNotification struct {
	Users []string
	Match *types.Match
}

type matchStartedNotification struct {
	*MatchStartedNotification
	cancel chan bool
}

type MatchEventNotifier struct {
	startedNotifications map[int16]*matchStartedNotification
	MatchStarted         chan *MatchStartedNotification
	mux                  sync.Mutex
	cancel               <-chan bool
}

func NewMatchEventNotifier(cancel <-chan bool) *MatchEventNotifier {
	matchEventNotifier := &MatchEventNotifier{
		startedNotifications: make(map[int16]*matchStartedNotification),
		MatchStarted:         make(chan *MatchStartedNotification),
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
		notification := &matchStartedNotification{
			MatchStartedNotification: &MatchStartedNotification{
				Users: []string{userID},
				Match: match,
			},
			cancel: make(chan bool, 1),
		}

		if r.startMatchTicker(notification) {
			r.startedNotifications[match.ID] = notification
		}
	}
}

func (r *MatchEventNotifier) RemoveUserFromNotificationsForMatch(match *types.Match, userID string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if existing, ok := r.startedNotifications[match.ID]; ok {
		existing.Users = slices.DeleteFunc(existing.Users, func(existingUserID string) bool {
			return existingUserID == userID
		})

		// If there's no users left, cancel the notification and remove it from the map
		if len(existing.Users) == 0 {
			log.Printf("No users left for match %d, closing cancel channel", match.ID)
			close(existing.cancel)
		}
	}
}

func (r *MatchEventNotifier) GetSubscribedMatchesForUser(userID string) []*types.Match {
	r.mux.Lock()
	defer r.mux.Unlock()

	var matches []*types.Match
	for _, value := range r.startedNotifications {
		if slices.Contains(value.Users, userID) {
			matches = append(matches, value.Match)
		}
	}
	return matches
}

func (r *MatchEventNotifier) startMatchTicker(notification *matchStartedNotification) bool {
	now := time.Now().UTC()
	matchStart := time.Unix(notification.Match.ScheduledTime, 0).UTC()
	duration := matchStart.Sub(now)

	if duration > 0 {
		go func() {
			ticker := time.NewTicker(duration)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					{

						r.mux.Lock()
						delete(r.startedNotifications, notification.Match.ID)
						r.mux.Unlock()

						r.MatchStarted <- notification.MatchStartedNotification
						return
					}
				case <-notification.cancel:
					{
						log.Printf("Notification being cancelled for match: %d", notification.Match.ID)
						return
					}
				case <-r.cancel:
					{
						return
					}
				}
			}
		}()
		return true
	}
	return false
}
