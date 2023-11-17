package bot

import (
	"github.com/flusaka/dota-tournament-bot/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSingleUserReceivesNotificationOfMatchStart(t *testing.T) {
	cancel := make(chan bool, 1)
	matchNotifier := NewMatchEventNotifier(cancel)
	match := &types.Match{
		ID:            0,
		TeamOne:       &types.Team{DisplayName: "OG"},
		TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
		ScheduledTime: time.Now().Add(time.Second * 1).UTC().Unix(),
		StreamUrl:     "https://twitch.tv",
	}
	uid := "userid1"
	matchNotifier.AddUserToNotificationsForMatch(match, uid)
	matchNotifier.AddUserToNotificationsForMatch(match, uid)

	assert.Contains(t, matchNotifier.startedNotifications, match.ID)

	notification := <-matchNotifier.MatchStarted
	assert.NotNil(t, notification)
	assert.NotNil(t, notification.Match)
	assert.Equal(t, notification.Match.ID, match.ID)
	assert.Len(t, notification.Users, 1)
	assert.Equal(t, notification.Users[0], uid)
	assert.NotContains(t, matchNotifier.startedNotifications, match.ID)
}

func TestMultipleUsersReceivesNotificationOfMatchStart(t *testing.T) {
	cancel := make(chan bool, 1)
	matchNotifier := NewMatchEventNotifier(cancel)
	users := []string{"userid1",
		"userid2",
		"userid3",
		"userid4",
		"userid5",
	}
	match := &types.Match{
		ID:            0,
		TeamOne:       &types.Team{DisplayName: "OG"},
		TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
		ScheduledTime: time.Now().Add(time.Second * 1).UTC().Unix(),
		StreamUrl:     "https://twitch.tv",
	}
	for _, uid := range users {
		matchNotifier.AddUserToNotificationsForMatch(match, uid)
	}

	assert.Contains(t, matchNotifier.startedNotifications, match.ID)

	notification := <-matchNotifier.MatchStarted
	assert.NotNil(t, notification)
	assert.NotNil(t, notification.Match)
	assert.Equal(t, notification.Match.ID, match.ID)
	assert.Len(t, notification.Users, len(users))
	for i, uid := range notification.Users {
		assert.Equal(t, uid, users[i])
	}

	assert.NotContains(t, matchNotifier.startedNotifications, match.ID)
}

func TestSingleUserReceivesNotificationsOfMultipleMatchStarts(t *testing.T) {
	cancel := make(chan bool, 1)
	matchNotifier := NewMatchEventNotifier(cancel)
	match := &types.Match{
		ID:            0,
		TeamOne:       &types.Team{DisplayName: "OG"},
		TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
		ScheduledTime: time.Now().Add(time.Second * 1).UTC().Unix(),
		StreamUrl:     "https://twitch.tv",
	}
	anotherMatch := &types.Match{
		ID:            1,
		TeamOne:       &types.Team{DisplayName: "Gaimin Gladiators"},
		TeamTwo:       &types.Team{DisplayName: "9Pandas"},
		ScheduledTime: time.Now().Add(time.Second * 2).UTC().Unix(),
		StreamUrl:     "https://twitch.tv",
	}
	uid := "userid1"
	matchNotifier.AddUserToNotificationsForMatch(match, uid)
	matchNotifier.AddUserToNotificationsForMatch(anotherMatch, uid)

	assert.Contains(t, matchNotifier.startedNotifications, match.ID)
	assert.Contains(t, matchNotifier.startedNotifications, anotherMatch.ID)

	notification := <-matchNotifier.MatchStarted
	assert.NotNil(t, notification)
	assert.NotNil(t, notification.Match)
	assert.Equal(t, notification.Match.ID, match.ID)
	assert.Len(t, notification.Users, 1)
	assert.Equal(t, notification.Users[0], uid)

	assert.NotContains(t, matchNotifier.startedNotifications, match.ID)

	notification = <-matchNotifier.MatchStarted
	assert.NotNil(t, notification)
	assert.NotNil(t, notification.Match)
	assert.Equal(t, notification.Match.ID, anotherMatch.ID)
	assert.Len(t, notification.Users, 1)
	assert.Equal(t, notification.Users[0], uid)

	assert.NotContains(t, matchNotifier.startedNotifications, anotherMatch.ID)
}

func TestMultipleUsersReceivesNotificationsOfMultipleMatchStarts(t *testing.T) {
	cancel := make(chan bool, 1)
	matchNotifier := NewMatchEventNotifier(cancel)
	users := []string{"userid1",
		"userid2",
		"userid3",
		"userid4",
		"userid5",
	}
	matches := []*types.Match{
		{
			ID:            0,
			TeamOne:       &types.Team{DisplayName: "OG"},
			TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
			ScheduledTime: time.Now().Add(time.Second * 1).UTC().Unix(),
			StreamUrl:     "https://twitch.tv",
		},
		{
			ID:            1,
			TeamOne:       &types.Team{DisplayName: "OG"},
			TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
			ScheduledTime: time.Now().Add(time.Second * 2).UTC().Unix(),
			StreamUrl:     "https://twitch.tv",
		}, {
			ID:            2,
			TeamOne:       &types.Team{DisplayName: "OG"},
			TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
			ScheduledTime: time.Now().Add(time.Second * 3).UTC().Unix(),
			StreamUrl:     "https://twitch.tv",
		},
		{
			ID:            3,
			TeamOne:       &types.Team{DisplayName: "OG"},
			TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
			ScheduledTime: time.Now().Add(time.Second * 4).UTC().Unix(),
			StreamUrl:     "https://twitch.tv",
		},
		{
			ID:            4,
			TeamOne:       &types.Team{DisplayName: "OG"},
			TeamTwo:       &types.Team{DisplayName: "Team Liquid"},
			ScheduledTime: time.Now().Add(time.Second * 5).UTC().Unix(),
			StreamUrl:     "https://twitch.tv",
		},
	}
	for i := range users {
		matchNotifier.AddUserToNotificationsForMatch(matches[i], users[i])
		assert.Contains(t, matchNotifier.startedNotifications, matches[i].ID)
	}

	index := 0
	for {
		select {
		case notification := <-matchNotifier.MatchStarted:
			{
				assert.NotNil(t, notification)
				assert.NotNil(t, notification.Match)
				assert.Equal(t, notification.Match.ID, matches[index].ID)
				assert.Len(t, notification.Users, 1)
				assert.Equal(t, notification.Users[0], users[index])
				assert.NotContains(t, matchNotifier.startedNotifications, matches[index].ID)
				index++
				if index >= len(users) {
					return
				}
			}
		}
	}
}
