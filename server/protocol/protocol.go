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

type LogEvent struct {
	LogType  int    `json:"logType"`
	ProtagID uint64 `json:"protagID"`
	AntagID  uint64 `json:"antagID"`
}

type WorldState struct {
	Time     uint64   `json:"time"`
	Entities []Entity `json:"entities"`
}

type Event struct {
	Type string
	Data interface{}
}
