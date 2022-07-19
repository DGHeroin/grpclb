package xclient

import (
    "context"
    . "github.com/DGHeroin/grpclb/client"
    handler2 "github.com/DGHeroin/grpclb/common/handler"
)

type (
    XClient struct {
        client       Client
        pushHandlers *handler2.PushHandlers
    }
)

func NewXClient(discover Discover) XClient {
    cli := XClient{
        pushHandlers: handler2.NewPushHandlers(),
    }
    cli.client = NewClient(discover, WithPushMessage(cli.onPushRaw))
    return cli
}
func (cli *XClient) RegisterPush(name string, fn handler2.PushHandlerFunc) bool {
    return cli.pushHandlers.Register(name, fn)
}
func (cli *XClient) onPush(ctx context.Context, name string, r, w interface{}) error {
    data, err := handler2.Marshal(r)
    if err != nil {
        return err
    }
    reply, err := cli.client.Send(ctx, name, data)
    if err != nil {
        return err
    }
    return handler2.Unmarshal(reply, w)
}
func (cli *XClient) Request(ctx context.Context, name string, r, w interface{}) error {
    data, err := handler2.Marshal(r)
    if err != nil {
        return err
    }
    reply, err := cli.client.Send(ctx, name, data)
    if err != nil {
        return err
    }
    return handler2.Unmarshal(reply, w)
}

func (cli *XClient) onPushRaw(name string, payload []byte) ([]byte, error) {
    return cli.pushHandlers.HandleMessage(nil, name, payload)
}
