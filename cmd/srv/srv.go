package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/DGHeroin/grpclb/server"
    "log"
    "net"
    "sync/atomic"
    "time"
)

var (
    address string
    qps     int64
)

func init() {
    flag.StringVar(&address, "addr", ":30001", "serve address")
    flag.Parse()

    go func() {
        for {
            time.Sleep(time.Second)
            last := atomic.LoadInt64(&qps)
            atomic.StoreInt64(&qps, 0)
            log.Println("qps", last)
        }
    }()
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

func (s srv) OnMessage(ctx context.Context, name string, payload []byte) ([]byte, error) {
    // log.Println("收到消息", name, string(payload))
    atomic.AddInt64(&qps, 1)
    return []byte(fmt.Sprintf("%s/%v", name, time.Now())), nil
}
