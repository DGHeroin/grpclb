package xserver

import (
    "context"
    "github.com/DGHeroin/grpclb/handler"
    . "github.com/DGHeroin/grpclb/server"
    "net"
)

type (
    XServer struct {
        impl *serverImpl
    }
    serverImpl struct {
        opt      *option
        server   Server
        handlers *handler.Handlers
    }
    option struct {
        onPushClientNew   func(client PushClient)
        OnPushClientClose func(client PushClient)
    }
    OptionFunc func(*option)
)

func (s *serverImpl) OnMessage(ctx context.Context, name string, payload []byte) ([]byte, error) {
    return s.handlers.HandleMessage(ctx, name, payload)
}

func (s *serverImpl) OnPushClientNew(client PushClient) {
    if s.opt.onPushClientNew != nil {
        s.opt.onPushClientNew(client)
    }
}

func (s *serverImpl) OnPushClientClose(client PushClient) {
    if s.opt.OnPushClientClose != nil {
        s.opt.OnPushClientClose(client)
    }
}

func NewXServer(opts ...OptionFunc) XServer {
    o := &option{}
    for _, opt := range opts {
        opt(o)
    }
    s := XServer{
        impl: &serverImpl{
            opt:      o,
            handlers: handler.NewHandlers(),
        },
    }
    s.impl.server = NewServer(s.impl)

    return s
}

func (s *XServer) RegisterPush(name string, fn interface{}) bool {
    return s.impl.handlers.Register(name, fn)
}

func (s *XServer) ServeListener(ln net.Listener) error {
    return s.impl.server.ServeListener(ln)
}

func WithPushEvent(onNew, onClose func(client PushClient)) OptionFunc {
    return func(o *option) {
        o.onPushClientNew = onNew
        o.OnPushClientClose = onClose
    }
}
