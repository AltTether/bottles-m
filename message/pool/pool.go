package pool


type Message interface {
	Text() *string
}

type Pool interface {
	Get() (Message, error)
	Post(Message) error
}
