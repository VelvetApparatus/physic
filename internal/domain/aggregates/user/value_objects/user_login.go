package value_objects

const (
	loginMinLen = 6
	loginMaxLen = 30
)

type UserLogin string

func (ul UserLogin) String() string { return string(ul) }

// todo: what about errors ?

func ParseLogin(login string) (res UserLogin, valid bool) {
	if len(login) < loginMinLen || len(login) > loginMaxLen {
		return "", false
	}

	return UserLogin(login), true
}
