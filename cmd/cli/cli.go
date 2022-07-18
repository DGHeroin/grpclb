package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/DGHeroin/grpclb/client"
    "log"
    "os"
    "strings"
    "sync/atomic"
    "time"
)

var (
    address = "localhost:30001,localhost:30002"
)

func init() {
    flag.StringVar(&address, "addr", "ts.1kb.win", "server address")
    flag.Parse()
}

type discovery struct {
    ch chan []string
}

func (d *discovery) WatchUpdate() <-chan []string {
    return d.ch
}

func main() {
    dis := &discovery{
        ch: make(chan []string),
    }
    go func() {
        dis.ch <- strings.Split(address, ",")
    }()
    go func() {
        r := bufio.NewReader(os.Stdin)
        for {
            line, _, err := r.ReadLine()
            if err == nil {
                str := strings.TrimSpace(string(line))
                if str == "" {
                    continue
                }
                infos := strings.Split(str, ",")
                log.Println(infos)
                dis.ch <- infos
            }
        }
    }()
    qps := int32(0)
    go func() {
        for {
            time.Sleep(time.Second)
            now := atomic.LoadInt32(&qps)
            atomic.StoreInt32(&qps, 0)
            if now == 0 {
                continue
            }
            fmt.Println(now)
        }
    }()
    cli := client.NewClient(dis, client.WithPushMessage(func(name string, payload []byte) ([]byte, error) {
        // log.Println("收到推送哎")
        atomic.AddInt32(&qps, 1)
        return nil, nil
    }))

    request(cli, time.Second)
}
func request(cli client.Client, duration time.Duration) {
    i := 0
    for {
        i++
        if duration > 0 {
            time.Sleep(duration)
        }
        _, err := cli.Send(nil, "hello", []byte(fmt.Sprintf("hello world:%d", i)))
        if err != nil {
            log.Println(err)
            continue
        }
    }
}
