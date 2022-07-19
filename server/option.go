package server

type (
    option struct {
        onPushClientNew   func(client PushClient)
        onPushClientClose func(client PushClient)
    }
    Option func(*option)
)

func (o *option) OnPushClientNew(ps *pushClientImpl) {
    if o.onPushClientNew == nil {
        return
    }
    o.onPushClientNew(ps)
}

func (o *option) OnPushClientClose(ps *pushClientImpl) {
    if o.onPushClientClose == nil {
        return
    }
    o.onPushClientClose(ps)
}

func defaultOption() *option {
    return &option{}
}
func WithPushClientNew(fn func(client PushClient)) Option {
    return func(o *option) {
        o.onPushClientNew = fn
    }
}

func WithPushClientClose(fn func(client PushClient)) Option {
    return func(o *option) {
        o.onPushClientClose = fn
    }
}
