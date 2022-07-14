package server

import (
    "context"
    "github.com/DGHeroin/grpclb/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "net"
)

type (
    Server interface {
        ServeListener(ln net.Listener) error
    }
    Handler interface {
        OnMessage(ctx context.Context, name string, payload []byte) ([]byte, error)
        OnPushClientNew(client PushClient)
        OnPushClientClose(client PushClient)
    }
    serverImpl struct {
        rpcServer *grpc.Server
        handler   Handler
    }
)

func (s *serverImpl) Push(server pb.MessageHandler_PushServer) error {
    ps := newPushClient(func(name string, payload []byte) error {
        req := &pb.Message{
            Name:    name,
            Payload: payload,
        }
        return server.Send(req)
    })
    s.handler.OnPushClientNew(ps)

    ctx := server.Context()
    defer func() {
        s.handler.OnPushClientClose(ps)
    }()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (s *serverImpl) Request(ctx context.Context, r *pb.Message) (*pb.Message, error) {
    data, err := s.handler.OnMessage(ctx, r.Name, r.Payload)
    resp := &pb.Message{
        Payload: data,
    }
    if err != nil {
        resp.Error = []byte(err.Error())
    }
    return resp, nil
}
func (s *serverImpl) ServeListener(ln net.Listener) error {
    return s.rpcServer.Serve(ln)
}
func NewServer(handler Handler) Server {
    srv := &serverImpl{
        handler: handler,
    }
    s := grpc.NewServer()
    pb.RegisterMessageHandlerServer(s, srv)
    reflection.Register(s)

    srv.rpcServer = s

    return srv
}
