package server

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/common/errs"
    "github.com/DGHeroin/grpclb/common/pb"
    "github.com/DGHeroin/grpclb/common/utils"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "log"
    "net"
    "sync"
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

func (s *serverImpl) RegisterPush(server pb.MessageHandler_RegisterPushServer) error {
    utils.Recover()
    type pushContext struct {
        wg           *sync.WaitGroup
        requestData  []byte
        responseData []byte
        err          error
        name         string
    }
    sendCh := make(chan *pushContext)

    ps := newPushClient(func(name string, payload []byte) ([]byte, error) {
        utils.Recover()
        ctx := &pushContext{
            name:        name,
            wg:          &sync.WaitGroup{},
            requestData: payload,
        }
        ctx.wg.Add(1)
        sendCh <- ctx
        ctx.wg.Wait()

        return ctx.responseData, ctx.err
    })
    s.handler.OnPushClientNew(ps)

    ctx := server.Context()
    defer func() {
        s.handler.OnPushClientClose(ps)
        close(sendCh)
    }()
    // 首个消息
    if msg, err := server.Recv(); err != nil {
        return err
    } else {
        if msg.Name != "ping" {
            return fmt.Errorf("first message should be 'ping'")
        }
        if err := server.Send(&pb.Message{Name: "pong"}); err != nil {
            return fmt.Errorf("first message response 'pong' error:%v", err)
        }
    }

    for {
        select {
        case ctx := <-sendCh:
            if ctx == nil {
                return nil // signal by close
            }
            // send push
            err := server.Send(&pb.Message{Name: ctx.name, Payload: ctx.requestData})
            if err != nil {
                log.Println("server send push error:", err)
                ctx.wg.Done()
                return nil
            }
            // wait response
            if pushResponseMessage, err := server.Recv(); err != nil {
                ctx.err = err
            } else {
                ctx.responseData = pushResponseMessage.Payload
                if pushResponseMessage.ErrorCode != 0 {
                    ctx.err = errs.GetError(pushResponseMessage.ErrorCode)
                }
            }
            ctx.wg.Done()
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
        resp.ErrorCode = errs.ErrCodeRequestHandlerInvoke
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
