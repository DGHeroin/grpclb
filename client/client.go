package client

import (
    "context"
    "fmt"
    "github.com/DGHeroin/grpclb/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/resolver"
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
)

func NewClient(discover Discover) Client {
    return &client{
        discover: discover,
    }
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
    msg := &pb.RequestMessage{
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
