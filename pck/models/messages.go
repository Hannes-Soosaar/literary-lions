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

// GetInstance returns the singleton instance of Message.
func GetInstance() *Message {
	once.Do(func() {
		instance = &Message{}
	})
	return instance
}

// SetSuccess sets the Success field of the singleton instance.
func (m *Message) SetSuccess(success string) {
	mu.Lock()
	defer mu.Unlock()
	m.Success = success
}

// GetSuccess returns the Success field of the singleton instance.
func (m *Message) GetSuccess() string {
	mu.Lock()
	defer mu.Unlock()
	return m.Success
}

// SetError sets the Error field of the singleton instance.
func (m *Message) SetError(err error) {
	mu.Lock()
	defer mu.Unlock()
	m.Error = err
}

// GetError returns the Error field of the singleton instance.
func (m *Message) GetError() error {
	mu.Lock()
	defer mu.Unlock()
	return m.Error
}