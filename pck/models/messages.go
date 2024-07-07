package models

import(

"sync"

)

type Message struct{
	Success string
	Error error
}

var (
	instance *Message
	once     sync.Once
	mu       sync.Mutex
)

func GetInstance() *Message {
	once.Do(func() {
		instance = &Message{}
	})
	return instance
}

func (m *Message) SetSuccess(success string) {
	mu.Lock()
	defer mu.Unlock()
	m.Success = success
}

func (m *Message) GetSuccess() string {
	mu.Lock()
	defer mu.Unlock()
	return m.Success
}

func (m *Message) SetError(err error) {
	mu.Lock()
	defer mu.Unlock()
	m.Error = err
}

func (m *Message) GetError() error {
	mu.Lock()
	defer mu.Unlock()
	return m.Error
}