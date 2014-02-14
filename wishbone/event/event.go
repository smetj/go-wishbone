package event

// import "fmt"

func NewEvent() Event {
	return Event{Header: make(map[string]map[string]interface{})}
}

type Event struct {
	Header map[string]map[string]interface{}
	Data   interface{}
}
