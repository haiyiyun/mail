package data

import (
	"net/http"
)

type Header interface {
	Date() string
	MessageId() string
	Sender() string
}

type Dataer interface {
	Header() http.Header
	Data() []byte
}

type Mailer interface {
	Header
	Dataer
}
