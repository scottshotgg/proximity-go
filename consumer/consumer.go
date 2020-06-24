package consumer

type (
	Consumer interface {
		ID() string

		Single() ([]byte, error)
		Stream(ch chan<- []byte)
	}
)
