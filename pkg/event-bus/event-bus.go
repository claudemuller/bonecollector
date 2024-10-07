package eventbus

const pollInitSize = 100

const (
	EventDead = iota
	EventDestroyEntity
	EventPlayerMove
)

type eventType int

type Event struct {
	Type eventType
	Args interface{}
}

type HandlerFn func(args ...interface{})

type EventBus struct {
	handlers map[eventType][]HandlerFn
	poll     []Event
}

func New() EventBus {
	return EventBus{
		handlers: make(map[eventType][]HandlerFn),
		poll:     make([]Event, pollInitSize),
	}
}

func (eb *EventBus) On(et eventType, handler HandlerFn) {
	eb.handlers[et] = append(eb.handlers[et], handler)
}

func (eb *EventBus) Emit(e Event) {
	eb.poll = append(eb.poll, e)
}

func (eb *EventBus) Process() {
	for _, e := range eb.poll {
		for _, h := range eb.handlers[e.Type] {
			h(e.Args)
		}
	}
	eb.poll = make([]Event, pollInitSize)
}
