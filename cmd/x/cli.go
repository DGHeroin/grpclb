package main

import (
    "context"
    "github.com/DGHeroin/grpclb/x/xclient"
    "log"
)

type discovery struct {
    ch chan []string
}

func (d *discovery) WatchUpdate() <-chan []string {
    return d.ch
}
func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    dis := &discovery{
        ch: make(chan []string, 2),
    }
    dis.ch <- []string{"localhost:30001"}
    cli := xclient.NewXClient(dis)

    cli.RegisterPush("on.client.push", func(ctx context.Context, payload []byte) ([]byte, error) {
        log.Println("收到消息", string(payload))
        return []byte("你好,推送消息!"), nil
    })

    var (
        r = "你好"
        w string
    )
    err := cli.Request(context.Background(), "say hello", &r, &w)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println(w)

    select {}
}
