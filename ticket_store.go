package cas

import (
	"errors"
)

var (
	ErrInvalidTicket = errors.New("cas: ticket store: invalid ticket")
)

type TicketStore interface {
	Read(id string) (*AuthenticationResponse, error)
	Write(id string, ticket *AuthenticationResponse) error
	Delete(id string) error
	Clear() error
}
