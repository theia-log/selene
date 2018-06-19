package model

type Event struct {
	ID        string
	Timestamp float64
	Source    string
	Tags      []string
	Content   string
}

func (ev *Event) Load(eventData string) (err error) {
	return
}

func (ev *Event) LoadBytes(eventData []byte) (err error) {
	return
}

func (ev *Event) Dump() (eventData string, err error) {
	return
}

func (ev *Event) DumpBytes() (eventData []byte, err error) {
	return
}
