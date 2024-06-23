package utils

type AppErrorType string

const (
	InvalidCharInMsgLength      AppErrorType = "Invalid char in message length"
	CannotConvertStringToNumber AppErrorType = "Cannot convert string to number"
	WrongCommandFormat          AppErrorType = "Wrong command format"
)

type AppError struct {
	ErrType AppErrorType
	Msg     string
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Type() AppErrorType {
	return e.ErrType
}

func GetErrorType(err error) AppErrorType {
	if re, ok := err.(*AppError); ok {
		return re.Type()
	}

	return ""
}
