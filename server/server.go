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
    }
    serverImpl struct {
        rpcServer *grpc.Server
        handler   Handler
    }
    Message struct {
    }
)

func (s *serverImpl) Request(ctx context.Context, r *pb.RequestMessage) (*pb.ResponseMessage, error) {
    data, err := s.handler.OnMessage(ctx, r.Name, r.Payload)
    resp := &pb.ResponseMessage{
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
