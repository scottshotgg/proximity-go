package consumer

type (
	Consumer interface {
		ID() string

		Recv() ([]byte, error)
		Listen(ch chan<- []byte)
	}
)
