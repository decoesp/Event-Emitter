package eventemitter

import "sync"

type listen func(data interface{})

type Eventmitter struct {
	mu     sync.RWMutex
	events map[string][]listen
	once   map[string][]listen
}

func StartMap() Eventmitter {
	return Eventmitter{
		events: make(map[string][]listen),
		once:   make(map[string][]listen),
	}

}

func (emitter *Eventmitter) Listen(eventName string, fn func(data any)) {

	emitter.mu.Lock()
	defer emitter.mu.Unlock()
	emitter.events[eventName] = append(emitter.events[eventName], fn)

}

func (emitter *Eventmitter) ListenOnce(eventName string, fn func(data any)) {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()
	emitter.once[eventName] = append(emitter.once[eventName], fn)

}

func (emitter *Eventmitter) Emit(eventName string, data any) {
	emitter.mu.Lock()

	listeners := append([]listen{}, emitter.events[eventName]...)
	onceListeners := append([]listen{}, emitter.once[eventName]...)

	delete(emitter.once, eventName)

	emitter.mu.Unlock()

	for _, fn := range listeners {
		safeCall(fn, data)
	}

	for _, fn := range onceListeners {
		safeCall(fn, data)
	}

}

func (emitter *Eventmitter) RemoveEvent(eventName string) {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()

	delete(emitter.events, eventName)
	delete(emitter.once, eventName)
}

func (emitter *Eventmitter) Reset() {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()

	emitter.events = make(map[string][]listen)
	emitter.once = make(map[string][]listen)
}

func safeCall(fn listen, data interface{}) {
	defer func() {
		if recover() != nil {

		}
	}()

	fn(data)
}
