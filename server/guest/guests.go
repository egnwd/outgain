package guest

const (
	UserType = iota
	BotType
)

type Guest struct {
	Type   int
	Name   string
	Source string
	gains  int
}

type List struct {
	List     []*Guest
	UserSize int
}

// NewUser returns a user with a specified name and no gains
func NewUser(name string, source string) *Guest {
	return &Guest{
		Type:   UserType,
		Name:   name,
		Source: source,
	}
}

// NewBot returns a bpt with a specified name
func NewBot(name string, source string) *Guest {
	return &Guest{
		Type:   BotType,
		Name:   name,
		Source: source,
	}
}

// GetName returns the name of the user
func (g *Guest) GetName() string {
	return g.Name
}

func (g *Guest) AddGains(amount int) {
	g.gains += amount
}

func (g *Guest) LoseGains(amount int) {
	g.gains -= amount
}

func (g *Guest) GetGains() int {
	return g.gains
}

func (g *Guest) ResetGains() {
	g.gains = 0
}

func (guests List) Iterator() []*Guest {
	return guests.List
}

func (g *Guest) ResetScore() {
	g.gains = 0
}

func (g *Guest) IsUser() bool {
	return g.Type == UserType
}
