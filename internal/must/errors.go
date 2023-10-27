package must

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
)

var (
	//common
	ErrInvalidCredentials  = &Error{Code: -9000, Message: "Invalid credentials."}
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

func handleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	if httpStatus != http.StatusMethodNotAllowed {
		runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
		return
	}

	// Use HTTPStatusError to customize the DefaultHTTPErrorHandler status code
	err := &runtime.HTTPStatusError{
		HTTPStatus: http.StatusInternalServerError,
		Err:        status.Errorf(http.StatusInternalServerError, ErrInternalServerError.Message),
	}

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
