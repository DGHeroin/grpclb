package client

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/pb"
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
    }
    Discover interface {
        WatchUpdate() <-chan []string
    }
    Client interface {
        Send(ctx context.Context, name string, payload []byte) ([]byte, error)
    }
    option struct {
        usePush  bool
        pushFunc func(name string, payload []byte) error
    }
    OptionFunc func(*option)
)

func NewClient(discover Discover, fns ...OptionFunc) Client {
    cli := &client{
        discover: discover,
    }
    opt := &option{}
    for _, fn := range fns {
        fn(opt)
    }
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

    ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
    conn, err := grpc.DialContext(ctx,
        fmt.Sprintf("%s://autority/%s", r.Scheme(), serviceName),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        c.onError(err)
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
    return resp.Payload, nil
}

func (c *client) onError(err error) {

}

func (c *client) loopPush() {
    for {
        c.waitPush()
    }
}
func (c *client) waitPush() {
    c.init()
    if c.conn == nil {
        return
    }
    cli := pb.NewMessageHandlerClient(c.conn)
    stream, err := cli.Push(context.Background())
    if err != nil {
        c.onError(err)
        return
    }
    req := &pb.Message{}
    err = stream.Send(req)
    if err != nil {
        log.Println("发送push失败:", err)
        return
    }
    for {
        resp, err := stream.Recv()
        if err != nil {
            log.Println("等待错误", err)
            break
        }
        log.Println("客户端收到推送:", resp.Name, string(resp.Payload))
    }
}
func WithPushMessage(fn func(name string, payload []byte) error) OptionFunc {
    return func(o *option) {
        o.usePush = true
        o.pushFunc = fn
    }
}
