package guest

const (
	UserType = iota
	BotType
)

type Guest struct {
	Type      int
	Name      string
	resources int
}

// NewUser returns a user with a specified name and no resources
func NewUser(name string) *Guest {
	return newGuest(name, UserType)
}

// NewBot returns a bpt with a specified name
func NewBot(name string) *Guest {
	return newGuest(name, BotType)
}

func newGuest(name string, t int) *Guest {
	return &Guest{Type: t, Name: name}
}

// GetName returns the name of the user
func (g *Guest) GetName() string {
	return g.Name
}

func (g *Guest) AddResources(amount int) {
	g.resources += amount
}

func (g *Guest) GetResources() int {
	return g.resources
}
