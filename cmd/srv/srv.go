package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/DGHeroin/grpclb/server"
    "log"
    "net"
    "time"
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
    s := server.NewServer(&srv{})
    err = s.ServeListener(ln)
    if err != nil {
        panic(err)
    }
}

type srv struct {
}

func (s srv) OnPushClientNew(client server.PushClient) {
    log.Println("新建推送客户端", client)
    time.AfterFunc(time.Second, func() {
        go func() {
            startTime := time.Now()
            for {
                err := client.Push("你好", []byte("消息"))
                if err != nil {
                    log.Println("推送失败", err)
                    break
                }
                if time.Now().Sub(startTime) > time.Second*10 {
                    break
                }
            }
        }()

    })
}

func (s srv) OnPushClientClose(client server.PushClient) {
    log.Println("关闭推送客户端", client)
}

func (s srv) OnMessage(ctx context.Context, name string, payload []byte) ([]byte, error) {
    log.Println("收到消息", name, string(payload))
    return []byte(fmt.Sprintf("%s/%v", name, time.Now())), nil
}
