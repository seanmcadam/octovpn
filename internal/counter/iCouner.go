package counter

type CounterWidth int

const CounterWidth32 CounterWidth = 32
const CounterWidth64 CounterWidth = 64

type Counter interface {
	ToByte() []byte
	Uint() interface{}
	Copy() Counter
	Width() CounterWidth
}

type CounterStruct interface {
	Next() Counter
	GetCountCh() <-chan Counter
	NewByteCounter([]byte) (Counter, error)
	Width() CounterWidth
}
