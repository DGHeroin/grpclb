package client

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/common/errs"
    "github.com/DGHeroin/grpclb/common/handler"
    "github.com/DGHeroin/grpclb/common/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/resolver"
    "log"
    "sync"
    "time"
)

type (
    client struct {
        conn         *grpc.ClientConn
        discover     Discover
        mu           sync.Mutex
        opt          *option
        pushHandlers *handler.PushHandlers
    }
    Discover interface {
        WatchUpdate() <-chan []string
        Latest() []string
    }
    Client interface {
        Request(ctx context.Context, name string, r, w interface{}) error
    }
)

func (c *client) Request(ctx context.Context, name string, r, w interface{}) error {
    var (
        data []byte
        err  error
    )
    if r != nil {
        data, err = handler.Marshal(r)
        if err != nil {
            return err
        }
    }
    reply, err := c.sendRaw(ctx, name, data)
    if err != nil {
        return err
    }
    if w == nil {
        return nil
    }
    return handler.Unmarshal(reply, w)
}

func NewClient(discover Discover, fns ...Option) Client {
    cli := &client{
        discover: discover,
    }
    opt := defaultOption()
    for _, fn := range fns {
        fn(opt)
    }
    cli.opt = opt

    return cli
}
func (c *client) init() {
    c.mu.Lock()
    if c.conn != nil {
        c.mu.Unlock()
        return
    }
    defer c.mu.Unlock()

    r := newResolver(c.discover)
    resolver.Register(r)

    ctx, cancel := context.WithTimeout(context.TODO(), c.opt.timeout)
    target := fmt.Sprintf("%s://authority/%s", r.Scheme(), c.opt.servicePath)
    conn, err := grpc.DialContext(ctx,
        target,
        grpc.WithResolvers(r),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    cancel()
    if err != nil {
        log.Println("Dial:", err)
        close(r.closeCh)
        return
    }
    c.conn = conn
    go c.loopPush()
}
func (c *client) sendRaw(ctx context.Context, name string, payload []byte) ([]byte, error) {
    if ctx == nil {
        ctx = context.TODO()
    }
    c.init()
    if c.conn == nil {
        return nil, fmt.Errorf("grpclb:client not connect")
    }
    cli := pb.NewMessageHandlerClient(c.conn)
    msg := &pb.Message{
        Name:    name,
        Payload: payload,
    }
    resp, err := cli.Request(ctx, msg)
    if err != nil {
        return nil, fmt.Errorf("grpclb Request:%v", err)
    }
    if resp.ErrorCode != 0 {
        return resp.Payload, errs.GetError(resp.ErrorCode)
    }
    return resp.Payload, nil
}

func (c *client) loopPush() {
    if c.opt.pushFunc == nil {
        return
    }
    c.waitPush()
    time.AfterFunc(time.Millisecond*10, func() {
        go c.loopPush()
    })
}
func (c *client) waitPush() {
    c.init()

    if c.conn == nil {
        return
    }
    cli := pb.NewMessageHandlerClient(c.conn)
    stream, err := cli.RegisterPush(context.Background())
    if err != nil {
        return
    }
    if err = stream.Send(&pb.Message{Name: "ping"}); err != nil {
        return
    }

    if msg, err := stream.Recv(); err != nil {
        return
    } else {
        if msg.Name != "pong" {
            return
        }
    }

    pushFunc := c.opt.pushFunc

    for {
        resp, err := stream.Recv()
        if err != nil {
            break
        }
        data, err := pushFunc(resp.Name, resp.Payload)
        respMsg := &pb.Message{
            Name:    resp.Name,
            Payload: data,
        }
        if err != nil {
            respMsg.ErrorCode = errs.ErrCodePushHandlerInvoke
        }
        if err := stream.Send(respMsg); err != nil {
            break
        }
    }
}
