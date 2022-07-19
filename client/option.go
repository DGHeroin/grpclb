package client

import "time"

type (
    option struct {
        pushFunc    func(name string, payload []byte) ([]byte, error)
        servicePath string
        schema      string
        timeout     time.Duration
    }
    Option func(*option)
)

func defaultOption() *option {
    return &option{
        schema:      "dns",
        servicePath: "default",
        timeout:     time.Second * 5,
    }
}
func WithPushMessage(fn func(name string, payload []byte) ([]byte, error)) Option {
    return func(o *option) {
        o.pushFunc = fn
    }
}
func WithTimeout(duration time.Duration) Option {
    return func(o *option) {
        o.timeout = duration
    }
}
