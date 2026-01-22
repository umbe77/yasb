package models

type Message struct{
	data map[string]any
}

func NewMessage() *Message{
	return &Message{
		data: make(map[string]any),
	}
}

func (m *Message) Get() map[string]any {
	return m.data
}