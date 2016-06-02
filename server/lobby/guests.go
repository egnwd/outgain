package lobby

import (
	"errors"
	"fmt"
	"log"
)

const (
	userType = iota
	botType  = iota
)

type guest struct {
	Type int
	name string
}

type guestList struct {
	list     []*guest
	userSize int
}

// User identified by unique GitHub name
type User struct {
	guest
	resources int
}

// Bot is a computer controlled AI
type Bot struct {
	guest
}

// NewUser returns a user with a specified name and no resources
func NewUser(name string) *User {
	return &User{guest: guest{Type: userType, name: name}}
}

// NewBot returns a bpt with a specified name
func NewBot(name string) *Bot {
	return &Bot{guest: guest{Type: botType, name: name}}
}

// GetName returns the name of the user
func (g *guest) GetName() string {
	return g.name
}

func generalPopulation(size int) guestList {
	var bots guestList

	for i := size; i > 0; i-- {
		name := fmt.Sprintf("Bot %d", i)
		bots.list = append(bots.list, &NewBot(name).guest)
	}

	return bots
}

func (lobby *Lobby) ContainsUser(name string) bool {
	for _, g := range lobby.Guests.list {
		if g.name == name {
			return true
		}
	}

	return false
}

// AddUser adds the specified user to the lobby, returning an error if the
// lobby is already at capacity, and running the engine if the user is
// the first to join
func (lobby *Lobby) AddUser(user *User) error {
	lobbyGuests := lobby.Guests.list

	// Check for bot to remove
	bot, newGuests := lobbyGuests[0], lobbyGuests[1:]
	if bot.Type != botType {
		return errors.New("Lobby full")
	}

	i := len(lobbyGuests) - (lobby.Guests.userSize + 1)
	newGuest := []*guest{&user.guest}
	newGuests = append(newGuests[:i], append(newGuest, newGuests[i:]...)...)
	lobby.Guests.userSize++

	log.Printf("%d\n", lobby.Guests.userSize)

	lobby.Guests.list = newGuests
	return nil
}

func (guests guestList) Iterator() []*guest {
	return guests.list
}
