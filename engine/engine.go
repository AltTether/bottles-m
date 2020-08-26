package engine


type Message struct {
	Text *string
}

type Token struct {
	Str *string
}

type Bottle struct {
	Message *Message
	Token   *Token
}
