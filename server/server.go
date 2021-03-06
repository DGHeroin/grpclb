package server

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/common/errs"
    "github.com/DGHeroin/grpclb/common/handler"
    "github.com/DGHeroin/grpclb/common/pb"
    "github.com/DGHeroin/grpclb/common/utils"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "log"
    "net"
    "sync"
    "sync/atomic"
)

type (
    Server interface {
        ServeListener(ln net.Listener) error
        BroadcastPush(name string, payload []byte)
        RegisterHandler(name string, fn interface{}) bool
    }
    serverImpl struct {
        option       *option
        rpcServer    *grpc.Server
        handlers     *handler.Handlers
        mu           sync.RWMutex
        pushClientId int32
        allClients   map[int32]*pushClientImpl
    }
)

func (s *serverImpl) RegisterHandler(name string, fn interface{}) bool {
    return s.handlers.Register(name, fn)
}

func (s *serverImpl) RegisterPush(server pb.MessageHandler_RegisterPushServer) error {
    utils.Recover()
    type pushContext struct {
        wg           *sync.WaitGroup
        requestData  []byte
        responseData []byte
        err          error
        name         string
    }
    sendCh := make(chan *pushContext, 10)

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

    ctx := server.Context()
    defer func() {
        s.removePushClient(ps)
        go s.option.OnPushClientClose(ps)
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

    s.addPushClient(ps)
    go s.option.OnPushClientNew(ps)

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
    data, err := s.handlers.HandleMessage(ctx, r.Name, r.Payload)
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
func (s *serverImpl) addPushClient(ps *pushClientImpl) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for {
        id := atomic.AddInt32(&s.pushClientId, 1)
        if _, ok := s.allClients[id]; !ok {
            s.allClients[id] = ps
            ps.id = id
            break
        }
    }
}
func (s *serverImpl) removePushClient(ps *pushClientImpl) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.allClients, ps.id)
}
func (s *serverImpl) BroadcastPush(name string, payload []byte) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, v := range s.allClients {
        _, _ = v.Push(name, payload)
    }
}
func NewServer(opts ...Option) Server {
    o := defaultOption()
    for _, opt := range opts {
        opt(o)
    }

    srv := &serverImpl{
        option:     o,
        handlers:   handler.NewHandlers(),
        allClients: map[int32]*pushClientImpl{},
    }
    s := grpc.NewServer()
    pb.RegisterMessageHandlerServer(s, srv)
    reflection.Register(s)

    srv.rpcServer = s

    return srv
}
