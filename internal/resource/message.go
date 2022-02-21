package resource

type Message struct {
	Result	bool
	Cuz		string
}

func newMessage() *Message {
	return &Message{true, "-"}
}

func (m *Message) Failed(c string) {
	m.Result = false
	m.Cuz = c
}

func (m *Message) Reset() {
	m.Result = true
	m.Cuz = "-"
}
