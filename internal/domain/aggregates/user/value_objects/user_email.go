package value_objects

import (
	"net/mail"
)

type UserEmail string

func (us UserEmail) String() string { return string(us) }

func ParseEmail(email string) (UserEmail, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}

	return UserEmail(addr.String()), true
}
