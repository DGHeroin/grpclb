package main

import (
    "context"
    "github.com/DGHeroin/grpclb/server"
    "github.com/DGHeroin/grpclb/x/xserver"
    "log"
    "net"
    "time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    s := xserver.NewXServer(xserver.WithPushEvent(func(client server.PushClient) {
        pushMsg := []byte("server say hello")

        time.AfterFunc(time.Second, func() {
            reply, err := client.Push("on.client.push", pushMsg)
            if err != nil {
                log.Println(err)
            }
            log.Println("推送 1:", string(reply))
            reply, err = client.Push("on.client.not_exist", pushMsg)
            if err != nil {
                log.Println(err)
            }
            log.Println("推送 2:", reply)
        })
    }, func(client server.PushClient) {

    }))
    s.RegisterPush("say hello", sayHello)

    ln, err := net.Listen("tcp", ":30001")
    err = s.ServeListener(ln)
    log.Println(err)
}

func sayHello(ctx context.Context, r, w *string) error {
    log.Println("收到消息", *r)
    *w = "hello world!"
    return nil
}
