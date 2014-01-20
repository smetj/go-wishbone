package event
// import "fmt"

func NewEvent() Event {
    return Event{Header:make(map[string]HeaderEntry)}
}
type Event struct {
    Header map[string]HeaderEntry
    Data interface{}
}
type HeaderEntry struct{
    TTL int
    Data string
}
func (h HeaderEntry) IncrementTTL(){
    h.TTL++
}
