package errs

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type CodeError interface {
	Code() int
	Msg() string
	Detail() string
	WithDetail(detail string) CodeError
	// Is 判断是否是某个错误, loose为false时, 只有错误码相同就认为是同一个错误, 默认为true
	Is(err error, loose ...bool) bool
	Wrap(msg ...string) error
	error
}
type codeError struct {
	code   int
	msg    string
	detail string
}

func WarpGrpcErr(err error) error {
	code_err := &codeError{}
	ok := errors.As(err, &code_err)
	if ok {
		return status.Error(codes.Code(code_err.code), code_err.CompleteMsg())
	}
	return status.Error(codes.Code(500), code_err.Error())
}
func (e *codeError) Code() int {
	return e.code
}

//func (e *codeError) GRPCErr() error {
//	return status.Error(codes.Code(e.code), e.CompleteMsg())
//}

func (e *codeError) Msg() string {
	return e.msg
}

func (e *codeError) Detail() string {
	return e.detail
}

func (e *codeError) WithDetail(detail string) CodeError {
	var d string
	if e.detail == "" {
		d = detail
	} else {
		d = e.detail + ", " + detail
	}
	return &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: d,
	}
}

func (e *codeError) Wrap(w ...string) error {
	return errors.Wrap(e, strings.Join(w, ", "))
}
func (e *codeError) CompleteMsg() string {
	if e.detail == "" {
		return e.msg
	}
	return e.msg + ":" + e.detail
}
func (e *codeError) Is(err error, loose ...bool) bool {
	if err == nil {
		return false
	}

	codeErr, ok := Unwrap(err).(CodeError)
	if ok {
		return codeErr.Code() == e.code
	}
	return false
}

func (e *codeError) Error() string {
	v := make([]string, 0, 3)
	v = append(v, strconv.Itoa(e.code), e.msg)
	if e.detail != "" {
		v = append(v, e.detail)
	}
	return strings.Join(v, " ")
}
func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return err
}

func Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}
	if len(msg) == 0 {
		return errors.WithStack(err)
	}
	return errors.Wrap(err, strings.Join(msg, ", "))
}
func NewCodeError(code int, msg string) CodeError {
	return &codeError{
		code: code,
		msg:  msg,
	}
}
