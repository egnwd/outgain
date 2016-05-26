package protocol

// client/src/protocol.ts should be kept in sync with this

type Entity struct {
	Id     uint64  `json:"id"`
	Name   *string `json:"name"`
	Color  string  `json:"color"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

type WorldState struct {
	Time     uint64   `json:"time"`
	Entities []Entity `json:"entities"`
}

type Event struct {
	Type string
	Data interface{}
}
