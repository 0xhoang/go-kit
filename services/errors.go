package services

var (
	ErrInternalServerError = &Error{Code: -9001, Message: "internal server error"}
	ErrInvalidPassword     = &Error{Code: -1001, Message: "invalid password."}
	ErrEmailNotExists      = &Error{Code: -1004, Message: "email doesn't exist."}
	ErrInactiveAccount     = &Error{Code: -1008, Message: "uour account is inactive."}
	ErrEmailIsNotVerified  = &Error{Code: -1016, Message: "email is not verified."}
	ErrInvalidArgument     = &Error{Code: -9000, Message: "invalid argument"}
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}
