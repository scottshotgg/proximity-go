package producer

type (
	Producer interface {
		Single(route string, contents []byte) error
		Stream(route string, ch <-chan []byte)
	}
)
