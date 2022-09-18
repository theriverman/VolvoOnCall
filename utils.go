package vocdriver

type Position struct {
	Longitude float64     `json:"longitude"`
	Latitude  float64     `json:"latitude"`
	Timestamp string      `json:"timestamp"`
	Speed     interface{} `json:"speed"`   // TODO: figure out the actual type
	Heading   interface{} `json:"heading"` // TODO: figure out the actual type
}
