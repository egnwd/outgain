package protocol

// client/src/protocol.ts should be kept in sync with this

type Entity struct {
	Id     uint64  `json:"id"`
	Name   *string `json:"name"`
	Color  string  `json:"color"`
	Sprite *string `json:"sprite"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

type WorldState struct {
	Time      uint64   `json:"time"`
	Entities  []Entity `json:"entities"`
	LogEvents []string `json:"logEvents"`
}

type WorldDiff struct {
	Time         uint64   `json:"time"`
	PreviousTime uint64   `json:"previousTime"`
	Modified     []Entity `json:"modified"`
	Removed      []uint64 `json:"removed"`
}

type Event struct {
	Type string
	Data interface{}
}
