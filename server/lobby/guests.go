package lobby

const (
	userType = iota
	botType  = iota
)

type guest struct {
	Type int
	name string
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
