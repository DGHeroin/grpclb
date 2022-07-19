package server

type (
    PushClient interface {
        Push(name string, msg []byte) ([]byte, error)
    }
    pushClientImpl struct {
        pushFn pushFunc
        id     int32
    }
    pushFunc func(name string, payload []byte) ([]byte, error)
)

func (p pushClientImpl) Push(name string, msg []byte) ([]byte, error) {
    return p.pushFn(name, msg)
}

func newPushClient(pushFn pushFunc) *pushClientImpl {
    ps := &pushClientImpl{
        pushFn: pushFn,
    }
    return ps
}
