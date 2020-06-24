package consumer

type (
	Consumer interface {
		Single() ([]byte, error)
		Stream(ch chan<- []byte)
	}
)
