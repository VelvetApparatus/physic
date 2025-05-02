package value_objects

const (
	usernameMinLength = 5
	usernameMaxLength = 35
)

type UserName string

func (un UserName) String() string { return string(un) }

func ParseUsername(username string) (UserName, bool) {
	if len(username) < usernameMinLength || len(username) > usernameMaxLength {
		return "", false
	}

	return UserName(username), true
}
