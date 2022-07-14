package server

type (
    PushClient interface {
        Push(name string, payload []byte) error
    }
    pushClientImpl struct {
        pushFn pushFunc
    }
    pushFunc func(name string, payload []byte) error
)

func (p pushClientImpl) Push(name string, payload []byte) error {
    return p.pushFn(name, payload)
}

func newPushClient(pushFn pushFunc) *pushClientImpl {
    ps := &pushClientImpl{
        pushFn: pushFn,
    }
    return ps
}
