package value_objects

const (
	User  UserRole = `user`
	Admin UserRole = `admin`
)

type UserRole string

func (ur UserRole) String() string { return string(ur) }

func ParseUserRole(role string) (UserRole, bool) {
	switch UserRole(role) {
	case Admin, User:
		return UserRole(role), true
	default:
		return "", false
	}
}

func (ur UserRole) IsAdmin() bool {
	return ur == Admin
}
