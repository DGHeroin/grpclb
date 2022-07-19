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
    address string
)

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    flag.StringVar(&address, "addr", "localhost:30001,localhost:30002", "server address")
    flag.Parse()
}

type discovery struct {
    ch     chan []string
    latest []string
}

func (d *discovery) Latest() []string {
    return d.latest
}

func (d *discovery) WatchUpdate() <-chan []string {
    return d.ch
}
func (d *discovery) Dispatch(address []string) {
    d.ch <- address
    d.latest = address
}

func main() {
    dis := &discovery{
        ch: make(chan []string, 1),
    }
    dis.Dispatch(strings.Split(address, ","))
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
    cli := client.NewClient(dis,
        client.WithTimeout(time.Second),
        client.WithPushMessage(func(name string, payload []byte) ([]byte, error) {
            atomic.AddInt32(&qps, 1)
            log.Println("收到推送", name, string(payload))
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
        req := "你好呀"
        err := cli.Request(nil, "hello", &req, nil)
        if err != nil {
            log.Println("发送请求失败:", err)
            continue
        }
    }
}
