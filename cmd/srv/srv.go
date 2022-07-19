package main

import (
    "context"
    "flag"
    "github.com/DGHeroin/grpclb/server"
    "log"
    "net"
)

var (
    address string
)

func init() {
    flag.StringVar(&address, "addr", ":30001", "serve address")
    flag.Parse()
}

func main() {
    ln, err := net.Listen("tcp", address)
    if err != nil {
        panic(err)
    }
    s := server.NewServer(server.WithPushClientNew(func(client server.PushClient) {
        log.Println("注册推送客户端", client)
        msg := "你好 push client"
        client.Push("on.client.push", []byte(msg))
    }))
    s.RegisterHandler("hello", OnHello)
    err = s.ServeListener(ln)
    if err != nil {
        panic(err)
    }
}

func OnHello(ctx context.Context, r, w *string) error {
    log.Println(*r)
    return nil
}
