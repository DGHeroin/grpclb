package client

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/common/errs"
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
        conn     *grpc.ClientConn
        discover Discover
        mu       sync.Mutex
        opt      *option
    }
    Discover interface {
        WatchUpdate() <-chan []string
    }
    Client interface {
        Send(ctx context.Context, name string, payload []byte) ([]byte, error)
    }
    option struct {
        usePush  bool
        pushFunc func(name string, payload []byte) ([]byte, error)
        timeout  time.Duration
    }
    OptionFunc func(*option)
)

func NewClient(discover Discover, fns ...OptionFunc) Client {
    cli := &client{
        discover: discover,
    }
    opt := &option{
        timeout: time.Second * 5,
    }
    for _, fn := range fns {
        fn(opt)
    }
    cli.opt = opt
    if opt.usePush {
        go cli.loopPush()
    }
    return cli
}
func (c *client) init() {
    c.mu.Lock()
    if c.conn != nil {
        c.mu.Unlock()
        return
    }
    defer c.mu.Unlock()
    serviceName := "my-service"
    r := newResolver(c.discover, serviceName)
    resolver.Register(r)

    ctx, cancel := context.WithTimeout(context.TODO(), c.opt.timeout)
    conn, err := grpc.DialContext(ctx,
        fmt.Sprintf("%s://autority/%s", r.Scheme(), serviceName),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        c.onError(err)
        log.Println("发送错误", err)
    }
    cancel()
    c.conn = conn

}
func (c *client) Send(ctx context.Context, name string, payload []byte) ([]byte, error) {
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
        return nil, fmt.Errorf("grpclb:%v", err)
    }
    if resp.ErrorCode != 0 {
        return resp.Payload, errs.GetError(resp.ErrorCode)
    }
    return resp.Payload, nil
}

func (c *client) onError(err error) {
    // fmt.Println(err)
}

func (c *client) loopPush() {
    for {
        c.waitPush()
        time.Sleep(time.Second)
    }
}
func (c *client) waitPush() {
    c.init()
    if c.conn == nil {
        return
    }
    cli := pb.NewMessageHandlerClient(c.conn)
    stream, err := cli.RegisterPush(context.Background())
    if err != nil {
        log.Println("发生错误...", err)
        c.onError(err)
        return
    }

    if err = stream.Send(&pb.Message{Name: "ping"}); err != nil {
        log.Println("send ping:", err)
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
    if pushFunc == nil {
        return
    }
    for {
        resp, err := stream.Recv()
        if err != nil {
            log.Println("等待错误", err)
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
            log.Println("ack error:", err)
            break
        }
    }
}
func WithPushMessage(fn func(name string, payload []byte) ([]byte, error)) OptionFunc {
    return func(o *option) {
        o.usePush = true
        o.pushFunc = fn
    }
}
