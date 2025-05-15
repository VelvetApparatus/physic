package errors

var (
	InvalidPasswordError   = newUserError("invalid password")
	UsernameNotUniqueError = newUserError("username is not unique")
	LoginNotUniqueError    = newUserError("login not unique")
	EmailNotUniqueError    = newUserError("email not unique")
	InvalidUserRoleError   = newUserError("invalid user role")
	InvalidUserLoginError  = newUserError("invalid user login")
	InvalidUserEmailError  = newUserError("invalid email")
	InvalidUsernameError   = newUserError("invalid username")
)

type UserError struct {
	msg string
}

func (u *UserError) Error() string {
	return u.msg
}

func newUserError(msg string) *UserError {
	return &UserError{msg: msg}
}
