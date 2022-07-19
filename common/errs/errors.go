package errs

import "fmt"

var (
    ErrCodeHandlerNotFound      = int32(404)
    ErrCodePushHandlerInvoke    = int32(505)
    ErrCodeRequestHandlerInvoke = int32(606)
)
var (
    ErrHandlerNotFound      = fmt.Errorf("handler not found")
    ErrPushHandlerInvoke    = fmt.Errorf("push handle invoke error")
    ErrRequestHandlerInvoke = fmt.Errorf("request handle invoke error")
)
var errMap = map[int32]error{}

func init() {
    errMap[ErrCodeHandlerNotFound] = ErrHandlerNotFound
    errMap[ErrCodePushHandlerInvoke] = ErrPushHandlerInvoke
    errMap[ErrCodeRequestHandlerInvoke] = ErrRequestHandlerInvoke
}

func GetError(code int32) error {
    return errMap[code]
}
func IsError(a, b error) bool {
    return a == b
}
