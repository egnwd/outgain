package protocol

// client/src/protocol.ts should be kept in sync with this

type Creature struct {
	Id    uint64  `json:"id"`
	Name  string  `json:"name"`
	Color string  `json:"color"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Angle float64 `json:"angle"`
}

type WorldUpdate struct {
	Time      uint64     `json:"time"`
	Creatures []Creature `json:"creatures"`
}
