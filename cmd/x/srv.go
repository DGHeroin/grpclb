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
    var s xserver.XServer
    s = xserver.NewXServer(xserver.WithPushEvent(func(client server.PushClient) {
        pushMsg := []byte("server say hello")

        time.AfterFunc(time.Second, func() {
            reply, err := client.Push("on.client.push", pushMsg)
            if err != nil {
                log.Println(err)
            }
            log.Println("消息:", string(reply))
            s.BroadcastPush("on.client.broadcast", []byte("广播消息"))
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
