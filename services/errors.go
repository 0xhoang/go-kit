package services

var (
	//common
	ErrInternalServerError = &Error{Code: -9001, Message: "internal server error"}

	//user
	ErrInvalidPassword    = &Error{Code: -1001, Message: "invalid password."}
	ErrEmailNotExists     = &Error{Code: -1002, Message: "email doesn't exist."}
	ErrInactiveAccount    = &Error{Code: -1003, Message: "uour account is inactive."}
	ErrEmailIsNotVerified = &Error{Code: -1004, Message: "email is not verified."}
	ErrInvalidEmail       = &Error{Code: -1005, Message: "invalid email."}

	//argument
	ErrInvalidArgument = &Error{Code: -5000, Message: "invalid argument"}
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}
